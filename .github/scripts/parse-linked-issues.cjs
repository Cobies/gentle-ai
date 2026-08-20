'use strict';

// Shared parser for issue references in PR bodies.
// Single seam used by every pr-check.yml gate that reads linked issues:
// both check-issue-reference and check-issue-approved consume this parse.

const ANY_REFERENCE_PATTERN = /\b(?:(closes|fixes|resolves)|(refs?|references))\s+#(\d+)\b/gi;

// An HTML comment runs from `<!--` to the next `-->`, or to the end of the
// body when unclosed. The unclosed case matches how GitHub renders Markdown
// (the HTML spec closes an open comment at end of input), so everything after
// an unclosed `<!--` is invisible to reviewers and must not count.
const HTML_COMMENT_PATTERN = /<!--[\s\S]*?(?:-->|$)/g;

// Removes HTML comments so only reviewer-visible text remains.
function stripHtmlComments(body) {
  return (body || '').replace(HTML_COMMENT_PATTERN, '');
}

// Parses visible (non-comment) PR body for closing and non-closing issue references.
// Returns structured lists of issue numbers preserving order of first appearance.
function parseIssueReferences(body) {
  const visible = stripHtmlComments(body);
  const closing = [];
  const nonClosing = [];
  const all = [];
  const seenAll = new Set();
  const seenClosing = new Set();
  const seenNonClosing = new Set();

  for (const match of visible.matchAll(ANY_REFERENCE_PATTERN)) {
    const isClosing = Boolean(match[1]);
    const num = parseInt(match[3], 10);
    if (!Number.isSafeInteger(num) || num <= 0) {
      continue;
    }
    if (isClosing) {
      if (!seenClosing.has(num)) {
        seenClosing.add(num);
        closing.push(num);
      }
    } else {
      if (!seenNonClosing.has(num)) {
        seenNonClosing.add(num);
        nonClosing.push(num);
      }
    }
    if (!seenAll.has(num)) {
      seenAll.add(num);
      all.push(num);
    }
  }

  return {
    closing,
    nonClosing,
    all,
  };
}

// Backwards-compatible helper returning all referenced issue numbers (closing and non-closing).
function parseLinkedIssues(body) {
  return parseIssueReferences(body).all;
}

module.exports = {
  parseIssueReferences,
  parseLinkedIssues,
  stripHtmlComments,
};
