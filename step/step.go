package step

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/cache"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
)

const (
	stepId = "save-carthage-cache"

	// Cache key template
	// Arch: guarantees unique cache per stack
	// Checksum wildcards: in order to detect nested folders too, such as src/Cartfile.resolved or ios/Cartfile.resolved
	// Cachefile: this is the special cache fingerprint created by the Bitrise `carthage` step.
	// Since this file is created on the fly and isn't committed to the repo, it won't be available at cache restore
	// time, therefore we create a separate checksum so that the restore step can partially match the key using the
	// first checksum.
	key = `{{ .Arch }}-carthage-cache-{{ checksum "**/Cartfile.resolved" }}-{{ checksum "**/Carthage/Cachefile" }}`
)

var paths = []string{
	// Prebuilt frameworks
	// The `Carthage/Checkouts` folder is not needed, `carthage bootstrap --cache-builds` can avoid a rebuild without it
	"**/Carthage/Build",

	// Special cache fingerprint created by the Bitrise `carthage` step to skip the bootstrap of possible
	"**/Carthage/Cachefile",
}

type Input struct {
	Verbose bool `env:"verbose,required"`
}

type SaveCacheStep struct {
	logger       log.Logger
	inputParser  stepconf.InputParser
	pathChecker  pathutil.PathChecker
	pathProvider pathutil.PathProvider
	pathModifier pathutil.PathModifier
	envRepo      env.Repository
}

func New(
	logger log.Logger,
	inputParser stepconf.InputParser,
	pathChecker pathutil.PathChecker,
	pathProvider pathutil.PathProvider,
	pathModifier pathutil.PathModifier,
	envRepo env.Repository,
) SaveCacheStep {
	return SaveCacheStep{
		logger:       logger,
		inputParser:  inputParser,
		pathChecker:  pathChecker,
		pathProvider: pathProvider,
		pathModifier: pathModifier,
		envRepo:      envRepo,
	}
}

func (step SaveCacheStep) Run() error {
	var input Input
	if err := step.inputParser.Parse(&input); err != nil {
		return fmt.Errorf("failed to parse inputs: %w", err)
	}
	stepconf.Print(input)
	step.logger.Println()
	step.logger.Printf("Cache key: %s", key)
	step.logger.Printf("Cache paths:")
	step.logger.Printf(strings.Join(paths, "\n"))
	step.logger.Println()

	step.logger.EnableDebugLog(input.Verbose)

	saver := cache.NewSaver(step.envRepo, step.logger, step.pathProvider, step.pathModifier, step.pathChecker)
	return saver.Save(cache.SaveCacheInput{
		StepId:      stepId,
		Verbose:     input.Verbose,
		Key:         key,
		Paths:       paths,
		IsKeyUnique: true,
	})
}
