# CI/CD Integration for Microservices Testing

This guide shows how to integrate the testing framework with your CI/CD pipeline.

## GitHub Actions Example

Create `.github/workflows/tests.yml`:

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.23']

    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}
    
    - name: Run unit tests
      run: make test-unit
    
    - name: Run all tests with race detector
      run: make test-race
    
    - name: Generate coverage report
      run: make test-coverage
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        files: ./coverage.out
        flags: unittests
        name: codecov-umbrella
```

## GitLab CI Example

Create `.gitlab-ci.yml`:

```yaml
stages:
  - test

test:unit:
  stage: test
  image: golang:1.23
  script:
    - make test-unit
  coverage: '/coverage: \d+\.\d+%/'

test:all:
  stage: test
  image: golang:1.23
  script:
    - make test-verbose
  artifacts:
    reports:
      coverage_report:
        coverage_format: cobertura
        path: coverage.out

test:race:
  stage: test
  image: golang:1.23
  script:
    - make test-race
```

## Pre-commit Hook

Add to `.git/hooks/pre-commit`:

```bash
#!/bin/bash
echo "Running pre-commit tests..."
make test-unit || exit 1

echo "Checking code with linters..."
make lint || exit 1

echo "✓ All checks passed"
exit 0
```

Make executable:
```bash
chmod +x .git/hooks/pre-commit
```

## Jenkins Pipeline

Create `Jenkinsfile`:

```groovy
pipeline {
    agent any
    
    stages {
        stage('Test') {
            steps {
                sh 'make test-unit'
            }
        }
        
        stage('Test with Race Detector') {
            steps {
                sh 'make test-race'
            }
        }
        
        stage('Coverage') {
            steps {
                sh 'make test-coverage'
                publishHTML([
                    reportDir: '.',
                    reportFiles: 'coverage.html',
                    reportName: 'Code Coverage'
                ])
            }
        }
    }
    
    post {
        always {
            junit '**/test-results.xml'
            publishCoverage adapters: [coberturaAdapter('coverage.out')]
        }
    }
}
```

## GitHub Actions with Coverage Gates

```yaml
name: Test Coverage

on: [push, pull_request]

jobs:
  coverage:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.23'
    
    - name: Run tests with coverage
      run: make test-coverage
    
    - name: Check coverage threshold
      run: |
        COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        MIN_COVERAGE=75
        if (( $(echo "$COVERAGE < $MIN_COVERAGE" | bc -l) )); then
          echo "Coverage ${COVERAGE}% is below minimum ${MIN_COVERAGE}%"
          exit 1
        fi
        echo "Coverage ${COVERAGE}% meets minimum ${MIN_COVERAGE}%"
```

## Docker Build with Tests

Create `Dockerfile`:

```dockerfile
FROM golang:1.23 as tester

WORKDIR /app
COPY . .

RUN go mod download
RUN make test-race

# Build stage
FROM golang:1.23 as builder

WORKDIR /app
COPY --from=tester /app .

RUN CGO_ENABLED=0 go build -o /app/bin/service ./cmd/core/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
COPY --from=builder /app/bin/service /usr/local/bin/

CMD ["service"]
```

## Makefile CI Target

Add to your Makefile:

```makefile
.PHONY: ci-test ci-all

ci-test:
	@echo "Running CI tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

ci-all: lint ci-test
	@echo "✓ CI checks passed"

ci-check-coverage:
	@coverage=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Coverage: $$coverage%"; \
	if [ "$$(echo $$coverage | cut -d. -f1)" -lt 75 ]; then \
		echo "Coverage below 75%"; \
		exit 1; \
	fi
```

Run in CI:
```bash
make ci-all
make ci-check-coverage
```

## Coverage Requirements Strategy

### Track Coverage Over Time

```bash
#!/bin/bash
# coverage-tracker.sh

DATE=$(date +%Y-%m-%d)
COVERAGE=$(go test -cover ./... | grep coverage | awk '{print $NF}' | sed 's/coverage://;s/%//')

echo "$DATE: $COVERAGE%" >> coverage_history.txt
```

### Set Minimum Coverage by Package

```makefile
check-coverage-threshold:
	@go test -coverprofile=coverage.out ./...
	@echo "Bill package:" && go tool cover -func=coverage.out | grep Bill | tail -1
	@echo "Messaging package:" && go tool cover -func=coverage.out | grep Messaging | tail -1
	@echo "Overall:" && go tool cover -func=coverage.out | grep total
```

## Service-Specific CI Checks

```yaml
# For each service independently
test-core:
  script:
    - make test-service SERVICE=core

test-crawler:
  script:
    - make test-service SERVICE=crawler

test-messaging:
  script:
    - make test-service SERVICE=messaging
```

## Matrix Testing

Test against multiple Go versions:

```yaml
test:
  strategy:
    matrix:
      go-version: ['1.21', '1.22', '1.23']
  steps:
    - uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}
    - run: make test-verbose
```

## Continuous Integration Checklist

Before each commit:
- [ ] `make test-unit` passes
- [ ] `make lint` passes (if configured)
- [ ] New tests added for new code

Before each pull request:
- [ ] `make test-verbose` passes
- [ ] `make test-race` passes
- [ ] Coverage report reviewed
- [ ] No data races detected

Before each release:
- [ ] `make test-coverage` shows >75% coverage
- [ ] All integration tests pass
- [ ] Performance tests pass
- [ ] No flaky tests

## Notification Setup

### Slack Notifications

```yaml
- name: Notify on failure
  if: failure()
  uses: 8398a7/action-slack@v3
  with:
    status: ${{ job.status }}
    text: 'Tests failed on ${{ github.ref }}'
    webhook_url: ${{ secrets.SLACK_WEBHOOK }}
```

### Email Notifications

Configure in your CI system's native notification settings.

## Performance Optimization

### Parallel Testing in CI

```bash
# Increase parallelism for faster CI
go test -parallel 16 -count=1 -race ./...
```

### Caching Dependencies

```yaml
- uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
    restore-keys: |
      ${{ runner.os }}-go-
```

## Monitoring & Alerts

### Track Flaky Tests

```bash
#!/bin/bash
# Run tests multiple times to find flaky tests
for i in {1..10}; do
    echo "Run $i:"
    go test -short ./... || echo "Failed on run $i"
done
```

### Create Dashboard

Use tools like:
- Codecov.io for coverage tracking
- GitHub Actions Dashboard
- Custom metrics dashboard

## Summary

Key CI/CD targets available via Makefile:

```
make test              # Fast unit tests
make test-verbose      # All tests + race detector  
make test-coverage     # Generate coverage report
make test-race         # Race detector only
make lint              # Code linting
```

Use these in your CI pipeline to ensure code quality and test coverage!
