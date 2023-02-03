### Examples

#### Minimal example
```yaml
steps:
- restore-carthage-cache@1: {}
- carthage@3:
    inputs:
      carthage_command: bootstrap
      carthage_options: --use-xcframeworks --platform iOS
- xcode-test@2:
    inputs:
      scheme: MyScheme
- save-carthage-cache@1: {}
```

Check out [Workflow Recipes](https://github.com/bitrise-io/workflow-recipes#-key-based-caching-beta) for other caching examples!
