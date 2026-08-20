'use strict';

const fs = require('fs');

const REVIEW_BUDGET_LIMIT = 400;

/**
 * Loads and validates the grandfather size exceptions JSON file.
 * Fails closed on any error (missing file, invalid JSON, unexpected schema).
 *
 * @param {string} filePath Path to grandfather exceptions JSON
 * @returns {Set<number>} Set of allowed grandfathered PR numbers
 */
function loadGrandfatherExceptions(filePath) {
  if (!filePath || typeof filePath !== 'string') {
    throw new Error('Grandfather exceptions file path must be a non-empty string');
  }

  if (!fs.existsSync(filePath)) {
    throw new Error(`Grandfather exceptions file not found at: ${filePath}`);
  }

  let raw;
  try {
    raw = fs.readFileSync(filePath, 'utf8');
  } catch (err) {
    throw new Error(`Failed to read grandfather exceptions file: ${err.message}`);
  }

  let data;
  try {
    data = JSON.parse(raw);
  } catch (err) {
    throw new Error(`Grandfather exceptions file is malformed JSON: ${err.message}`);
  }

  if (!data || typeof data !== 'object' || data.version !== 1 || !Array.isArray(data.allowed_prs)) {
    throw new Error('Grandfather exceptions file has invalid schema (expected version: 1, allowed_prs: array)');
  }

  for (const pr of data.allowed_prs) {
    if (typeof pr !== 'number' || !Number.isInteger(pr) || pr <= 0) {
      throw new Error(`Invalid PR number in grandfather exceptions: ${JSON.stringify(pr)}`);
    }
  }

  return new Set(data.allowed_prs);
}

/**
 * Evaluates whether a pull request complies with the 400-line review budget.
 *
 * @param {Object} params
 * @param {number} [params.additions=0]
 * @param {number} [params.deletions=0]
 * @param {Array<string|Object>} [params.labels=[]]
 * @param {number|string} [params.prNumber]
 * @param {Set<number>|Array<number>} [params.grandfatherSet=new Set()]
 * @returns {{ pass: boolean, passed: boolean, warning: boolean, total: number, limit: number, message: string }}
 */
function evaluatePrSize({ additions = 0, deletions = 0, labels = [], prNumber, grandfatherSet = new Set() } = {}) {
  const total = (Number(additions) || 0) + (Number(deletions) || 0);
  const limit = REVIEW_BUDGET_LIMIT;

  const labelList = Array.isArray(labels)
    ? labels.map((l) => (typeof l === 'string' ? l : l?.name || ''))
    : [];
  const hasException = labelList.includes('size:exception');

  const num = Number(prNumber);
  const isGrandfathered = grandfatherSet instanceof Set
    ? grandfatherSet.has(num)
    : Array.isArray(grandfatherSet)
      ? grandfatherSet.includes(num)
      : false;

  if (total <= limit) {
    return {
      pass: true,
      passed: true,
      warning: false,
      total,
      limit,
      message: `PR is within the ${limit}-line cognitive review budget (+${additions}/-${deletions} = ${total}).`,
    };
  }

  if (isGrandfathered && hasException) {
    return {
      pass: true,
      passed: true,
      warning: true,
      total,
      limit,
      message:
        `PR changes ${total} lines (+${additions}/-${deletions}), exceeding the ${limit}-line cognitive review budget.\n\n` +
        `⚠️ Grandfathered PR #${num} with size:exception allowed with warning.`,
    };
  }

  if (isGrandfathered && !hasException) {
    return {
      pass: false,
      passed: false,
      warning: false,
      total,
      limit,
      message:
        `PR changes ${total} lines (+${additions}/-${deletions}), exceeding the ${limit}-line cognitive review budget.\n\n` +
        `PR #${num} is grandfathered but requires the 'size:exception' label to proceed.`,
    };
  }

  return {
    pass: false,
    passed: false,
    warning: false,
    total,
    limit,
    message:
      `PR changes ${total} lines (+${additions}/-${deletions}), exceeding the ${limit}-line cognitive review budget.\n\n` +
      'Why this exists: reviews should stay small enough to complete in ~60 minutes without reviewer fatigue.\n' +
      'New oversized PRs are not permitted even with size:exception. Please split this work into chained or stacked PRs.',
  };
}

module.exports = {
  evaluatePrSize,
  loadGrandfatherExceptions,
  REVIEW_BUDGET_LIMIT,
};
