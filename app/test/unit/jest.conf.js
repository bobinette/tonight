const path = require('path');

module.exports = {
  rootDir: path.resolve(__dirname, '../../'),
  moduleFileExtensions: ['js', 'json', 'vue'],
  moduleNameMapper: {
    '\\.(jpg|jpeg|png|gif|eot|otf|webp|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga)$':
      '<rootDir>/__mocks__/fileMock.js',
    '^@/(.*)$': '<rootDir>/src/$1',
  },
  transform: {
    '^.+\\.js$': '<rootDir>/node_modules/babel-jest',
    '.*\\.(vue)$': '<rootDir>/node_modules/vue-jest',
  },
  testPathIgnorePatterns: ['<rootDir>/test/e2e'],
  snapshotSerializers: ['<rootDir>/node_modules/jest-serializer-vue'],
  setupFiles: ['<rootDir>/test/unit/setup'],
  testResultsProcessor: 'jest-junit',
  // mapCoverage: true,
  // coverageDirectory: '<rootDir>/test/unit/coverage',
  // collectCoverageFrom: [
  //   'src/**/*.{js,vue}',
  //   '!src/main.js',
  //   '!**/node_modules/**'
  // ]
};
