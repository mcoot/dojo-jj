# Code patterns

Coding practices in this repository.

## Dependency injection

We use a Dependency Injection (DI) pattern.

Business logic must be under internal/service, in Service structs which take
their dependencies via a constructor, which are wired up via a factory.

Dependencies in internal/dependencies are boundaries of the application and must
contain as little logic as possible.

## Unit testing

Unit tests are integration-style and mock only the external boundary of the
application (i.e. the dependencies from internal/dependencies). They must invoke
business logic from the relevant entrypoint and assert against the result and 
observable external behaviour, never internal state.

Tests must cover the entire behaviour as we develop with red-green-refactor TDD.

### Fakes

Boundary dependencies have fakes under internal/dependencies/mocks. These are
behavioural, state-based fakes that should maximise realistic behaviour. They 
should be kept as simple as possible. If a fake absolutely requires more complex
internal logic, it should have its own unit tests.