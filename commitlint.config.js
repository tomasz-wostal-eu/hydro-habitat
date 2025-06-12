module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [2, 'always', [
      'feat', 'fix', 'docs', 'style', 'refactor', 
      'perf', 'test', 'chore', 'ci', 'build', 'revert'
    ]],
    'references-empty': [2, 'never'],
    'subject-case': [2, 'never', ['upper-case', 'pascal-case']],
    'subject-max-length': [2, 'always', 100],
    'body-leading-blank': [2, 'always'],
    'footer-leading-blank': [2, 'always']
  },
  parserPreset: {
    parserOpts: {
      referenceActions: ['closes', 'fixes', 'resolves', 'refs'],
      issuePrefixes: ['#', 'PR-']
    }
  }
};
