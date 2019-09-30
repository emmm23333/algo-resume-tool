package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.quvideo.com/algo/algo-resume-tool/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("rootCmd execute err: ", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initAPP)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.server.json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initAPP() {
	initConfig()
	initLog()
}

func initLog() {
	fmt.Println("initializing logger...")
	level := viper.GetString("log.level")
	logLevel := zap.InfoLevel
	switch level {
	case "error":
		logLevel = zap.ErrorLevel
	case "debug":
		logLevel = zap.DebugLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "Error":
		logLevel = zap.ErrorLevel
	}
	common.Outer = &lumberjack.Logger{
		Filename:  viper.GetString("log.path"),
		MaxSize:   viper.GetInt("log.max"),    // megabytes
		MaxAge:    viper.GetInt("log.maxAge"), // days
		LocalTime: viper.GetBool("log.localtime"),
	}
	w := zapcore.AddSync(common.Outer)
	var m zapcore.WriteSyncer
	if viper.GetBool("log.stdout") {
		m = zapcore.NewMultiWriteSyncer(w, zapcore.AddSync(os.Stdout))
	} else {
		m = zapcore.NewMultiWriteSyncer(w)
	}

	var core zapcore.Core
	format := viper.GetString("log.format")
	switch format {
	case "json":
		fallthrough
	case "":
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(newLogEncoderConfig()),
			m,
			logLevel,
		)
	case "console":
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(newLogEncoderConfig()),
			m,
			logLevel,
		)
	default:
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(newLogEncoderConfig()),
			m,
			logLevel,
		)
	}

	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync()
	common.Log = logger.Sugar()
}

func newLogEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	fmt.Println("initializing config...")
	viper.SetConfigType("json")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName(".conf")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("using local config err: ", err.Error())
		os.Exit(1)
	}
}
