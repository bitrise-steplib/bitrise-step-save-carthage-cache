package step

import (
	"fmt"

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
	// Cachefile: this is the special cache fingerprint created by the Bitrise `carthage-install` step.
	// Since this file is created on the fly and isn't committed to the repo, it won't be available at cache restore
	// time, therefore we create a separate checksum so that the restore step can partially match the key using the
	// first checksum.
	key = `{{ .Arch }}-carthage-cache-{{ checksum "**/Cartfile.resolved" }}-{{ checksum "**/Carthage/Cachefile" }}`

	// Cached path (looking at nested folders too)
	path = "**/Carthage"
)

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
	step.logger.Printf(path)
	step.logger.Println()

	step.logger.EnableDebugLog(input.Verbose)

	saver := cache.NewSaver(step.envRepo, step.logger, step.pathProvider, step.pathModifier, step.pathChecker)
	return saver.Save(cache.SaveCacheInput{
		StepId:      stepId,
		Verbose:     input.Verbose,
		Key:         key,
		Paths:       []string{path},
		IsKeyUnique: true,
	})
}
