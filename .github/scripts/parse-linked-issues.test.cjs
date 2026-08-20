'use strict';

const { test } = require('node:test');
const assert = require('node:assert/strict');

const { parseLinkedIssues, parseIssueReferences } = require('./parse-linked-issues.cjs');

test('a closing reference that only appears inside an HTML comment is not a linked issue', () => {
  const body = [
    '## Summary',
    '',
    '<!--',
    'Link the approved issue this PR resolves, e.g.:',
    'Closes #42',
    '-->',
    '',
    'Some visible description without any reference.',
  ].join('\n');

  assert.deepEqual(parseLinkedIssues(body), []);
  assert.deepEqual(parseIssueReferences(body), { closing: [], nonClosing: [], all: [] });
});

test('a non-closing reference that only appears inside an HTML comment is ignored', () => {
  const body = [
    '## Summary',
    '',
    '<!--',
    'Refs #3154',
    '-->',
    '',
    'Some visible description without any reference.',
  ].join('\n');

  assert.deepEqual(parseLinkedIssues(body), []);
  assert.deepEqual(parseIssueReferences(body), { closing: [], nonClosing: [], all: [] });
});

test('a real reference outside a comment plus the template example inside one finds exactly the real one', () => {
  const body = [
    'Closes #1770',
    '',
    '<!--',
    'Example: Closes #42',
    '-->',
  ].join('\n');

  assert.deepEqual(parseLinkedIssues(body), [1770]);
  assert.deepEqual(parseIssueReferences(body), { closing: [1770], nonClosing: [], all: [1770] });
});

test('a single-line HTML comment is also ignored', () => {
  const body = 'Fixes #7\n<!-- Closes #42 --> trailing visible text';

  assert.deepEqual(parseLinkedIssues(body), [7]);
  assert.deepEqual(parseIssueReferences(body), { closing: [7], nonClosing: [], all: [7] });
});

test('multiple visible closing references are all preserved, in order', () => {
  const body = 'Closes #10\nFixes #11\nResolves #12';

  assert.deepEqual(parseLinkedIssues(body), [10, 11, 12]);
  assert.deepEqual(parseIssueReferences(body), {
    closing: [10, 11, 12],
    nonClosing: [],
    all: [10, 11, 12],
  });
});

test('non-closing references (Refs #N, Ref #N, References #N) are distinguished from closing references', () => {
  const body = [
    '## Summary',
    '',
    'Closes #3352',
    'Refs #3154',
    'Ref #2000',
    'References #2001',
  ].join('\n');

  assert.deepEqual(parseLinkedIssues(body), [3352, 3154, 2000, 2001]);
  assert.deepEqual(parseIssueReferences(body), {
    closing: [3352],
    nonClosing: [3154, 2000, 2001],
    all: [3352, 3154, 2000, 2001],
  });
});

test('non-closing reference only satisfies parse and pr-check without closing keywords', () => {
  const body = 'Intermediate epic PR. Refs #3154';

  assert.deepEqual(parseLinkedIssues(body), [3154]);
  assert.deepEqual(parseIssueReferences(body), {
    closing: [],
    nonClosing: [3154],
    all: [3154],
  });
});

test('duplicate references in the body are deduplicated while preserving first-appearance order', () => {
  const body = [
    'Closes #10',
    'Refs #20',
    'Closes #10',
    'Refs #20',
    'Fixes #30',
  ].join('\n');

  assert.deepEqual(parseLinkedIssues(body), [10, 20, 30]);
  assert.deepEqual(parseIssueReferences(body), {
    closing: [10, 30],
    nonClosing: [20],
    all: [10, 20, 30],
  });
});

test('case insensitivity for keywords', () => {
  const body = 'closes #10 fixes #11 RESOLVES #12 REFS #13 ref #14';

  assert.deepEqual(parseLinkedIssues(body), [10, 11, 12, 13, 14]);
  assert.deepEqual(parseIssueReferences(body), {
    closing: [10, 11, 12],
    nonClosing: [13, 14],
    all: [10, 11, 12, 13, 14],
  });
});

// Documented behavior: an unclosed `<!--` hides everything through the end of
// the body. This matches how GitHub renders Markdown (the HTML spec closes an
// open comment at end of input), so a reference after an unclosed `<!--` is
// invisible to reviewers and must not count as a linked issue.
test('an unclosed HTML comment hides the rest of the body, matching rendered output', () => {
  const body = 'Closes #1770\n<!-- forgot to close this comment\nCloses #42\nRefs #99';

  assert.deepEqual(parseLinkedIssues(body), [1770]);
  assert.deepEqual(parseIssueReferences(body), { closing: [1770], nonClosing: [], all: [1770] });
});

test('an empty or missing body yields no linked issues', () => {
  assert.deepEqual(parseLinkedIssues(''), []);
  assert.deepEqual(parseLinkedIssues(null), []);
  assert.deepEqual(parseLinkedIssues(undefined), []);
  assert.deepEqual(parseIssueReferences(''), { closing: [], nonClosing: [], all: [] });
  assert.deepEqual(parseIssueReferences(null), { closing: [], nonClosing: [], all: [] });
  assert.deepEqual(parseIssueReferences(undefined), { closing: [], nonClosing: [], all: [] });
});
