# logger

## a wrapper for slog

## Features:
1. LevelTrace
2. LevelStats
3. logger.Trace()
4. logger.Stats()


## Sample Code:
```
logOptions := NewLogOptions(slog.LevelInfo, true)

logFile, err := os.OpenFile("log/file.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
if err != nil {
    panic(fmt.Errorf("error opening log file: %s", err))
}

appLogger := NewLogger("JSON", logFile, logOptions)

slog.SetDefault(appLogger.Logger)

appLogger.Trace("message", "time", time.Now().Unix().UTC())
appLogger.Stats("message", "time", time.Now().Unix().UTC())
appLogger.Debug("message", "time", time.Now().Unix().UTC())
appLogger.Info("message", "time", time.Now().Unix().UTC())
appLogger.Warn("message", "time", time.Now().Unix().UTC())
appLogger.Error("message", "time", time.Now().Unix().UTC())

```