'use strict';

const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const os = require('node:os');

const {
  evaluatePrSize,
  loadGrandfatherExceptions,
  REVIEW_BUDGET_LIMIT,
} = require('./check-pr-size.cjs');

test('standard PR within 400-line budget passes', () => {
  const result = evaluatePrSize({ additions: 150, deletions: 50, labels: [], prNumber: 100 });
  assert.equal(result.pass, true);
  assert.equal(result.passed, true);
  assert.equal(result.warning, false);
  assert.equal(result.total, 200);
});

test('PR exactly at 400 lines budget passes', () => {
  const result = evaluatePrSize({ additions: 300, deletions: 100, labels: [], prNumber: 101 });
  assert.equal(result.pass, true);
  assert.equal(result.warning, false);
  assert.equal(result.total, 400);
});

test('oversized PR (401 lines) without grandfathering fails even with size:exception label', () => {
  const result = evaluatePrSize({
    additions: 301,
    deletions: 100,
    labels: ['size:exception'],
    prNumber: 999,
    grandfatherSet: new Set([100, 200]),
  });
  assert.equal(result.pass, false);
  assert.equal(result.passed, false);
  assert.equal(result.warning, false);
  assert.equal(result.total, 401);
  assert.match(result.message, /New oversized PRs are not permitted even with size:exception/);
});

test('oversized PR without grandfathering and without label fails', () => {
  const result = evaluatePrSize({ additions: 500, deletions: 200, labels: [], prNumber: 999 });
  assert.equal(result.pass, false);
  assert.equal(result.warning, false);
  assert.equal(result.total, 700);
});

test('oversized PR with grandfathering AND size:exception passes with warning', () => {
  const result = evaluatePrSize({
    additions: 450,
    deletions: 50,
    labels: ['size:exception'],
    prNumber: 100,
    grandfatherSet: new Set([100]),
  });
  assert.equal(result.pass, true);
  assert.equal(result.passed, true);
  assert.equal(result.warning, true);
  assert.equal(result.total, 500);
  assert.match(result.message, /Grandfathered PR #100 with size:exception/);
});

test('oversized PR with grandfathering supports label objects', () => {
  const result = evaluatePrSize({
    additions: 450,
    deletions: 50,
    labels: [{ name: 'size:exception' }, { name: 'type:feature' }],
    prNumber: 100,
    grandfatherSet: new Set([100]),
  });
  assert.equal(result.pass, true);
  assert.equal(result.warning, true);
});

test('oversized PR with grandfathering but WITHOUT size:exception fails', () => {
  const result = evaluatePrSize({
    additions: 450,
    deletions: 50,
    labels: ['type:feature'],
    prNumber: 100,
    grandfatherSet: new Set([100]),
  });
  assert.equal(result.pass, false);
  assert.equal(result.passed, false);
  assert.equal(result.warning, false);
  assert.equal(result.total, 500);
  assert.match(result.message, /requires the 'size:exception' label/);
});

test('loadGrandfatherExceptions parses valid grandfather json file', () => {
  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'pr-size-test-'));
  const filePath = path.join(tmpDir, 'valid.json');
  fs.writeFileSync(filePath, JSON.stringify({ version: 1, allowed_prs: [42, 100, 200] }));

  const set = loadGrandfatherExceptions(filePath);
  assert.equal(set instanceof Set, true);
  assert.equal(set.size, 3);
  assert.equal(set.has(42), true);
  assert.equal(set.has(100), true);
  assert.equal(set.has(200), true);
  assert.equal(set.has(300), false);

  fs.rmSync(tmpDir, { recursive: true, force: true });
});

test('loadGrandfatherExceptions loads repository grandfather file', () => {
  const repoFilePath = path.join(__dirname, '..', 'grandfather-size-exceptions.json');
  const set = loadGrandfatherExceptions(repoFilePath);
  assert.equal(set instanceof Set, true);
});

test('loadGrandfatherExceptions fails closed if file is missing', () => {
  assert.throws(
    () => loadGrandfatherExceptions('/non/existent/path/exceptions.json'),
    /Grandfather exceptions file not found/
  );
});

test('loadGrandfatherExceptions fails closed if path is invalid', () => {
  assert.throws(() => loadGrandfatherExceptions(''), /non-empty string/);
  assert.throws(() => loadGrandfatherExceptions(null), /non-empty string/);
});

test('loadGrandfatherExceptions fails closed if JSON is malformed', () => {
  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'pr-size-test-'));
  const filePath = path.join(tmpDir, 'malformed.json');
  fs.writeFileSync(filePath, '{ invalid json');

  assert.throws(() => loadGrandfatherExceptions(filePath), /malformed JSON/);
  fs.rmSync(tmpDir, { recursive: true, force: true });
});

test('loadGrandfatherExceptions fails closed if schema/version is invalid', () => {
  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'pr-size-test-'));

  const badVersion = path.join(tmpDir, 'bad-version.json');
  fs.writeFileSync(badVersion, JSON.stringify({ version: 2, allowed_prs: [] }));
  assert.throws(() => loadGrandfatherExceptions(badVersion), /invalid schema/);

  const missingAllowed = path.join(tmpDir, 'missing-allowed.json');
  fs.writeFileSync(missingAllowed, JSON.stringify({ version: 1 }));
  assert.throws(() => loadGrandfatherExceptions(missingAllowed), /invalid schema/);

  const badAllowedType = path.join(tmpDir, 'bad-allowed-type.json');
  fs.writeFileSync(badAllowedType, JSON.stringify({ version: 1, allowed_prs: "42" }));
  assert.throws(() => loadGrandfatherExceptions(badAllowedType), /invalid schema/);

  const invalidPrElement = path.join(tmpDir, 'invalid-pr-element.json');
  fs.writeFileSync(invalidPrElement, JSON.stringify({ version: 1, allowed_prs: [42, "100"] }));
  assert.throws(() => loadGrandfatherExceptions(invalidPrElement), /Invalid PR number/);

  const negativePrElement = path.join(tmpDir, 'negative-pr-element.json');
  fs.writeFileSync(negativePrElement, JSON.stringify({ version: 1, allowed_prs: [-5] }));
  assert.throws(() => loadGrandfatherExceptions(negativePrElement), /Invalid PR number/);

  fs.rmSync(tmpDir, { recursive: true, force: true });
});
