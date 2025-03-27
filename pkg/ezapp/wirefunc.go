package ezapp

// WireFunc is a function type that wires together application components.
//
// A WireFunc takes a configuration of type C and returns a WireBundle containing
// the Runnables and CleanupFunc for the application. It is responsible for:
// 1. Creating and configuring application components using the provided configuration
// 2. Assembling the components into Runnables
// 3. Providing a cleanup function to release resources
//
// The generic type parameter C represents the configuration struct type.
// This is the same type that will be populated with environment variables by the Build function.
//
// Example:
//
//	func wireApp(cfg Config) (ezapp.WireBundle, error) {
//		// Create and configure components
//		db, err := database.Connect(cfg.DatabaseURL)
//		if err != nil {
//			return ezapp.WireBundle{}, fmt.Errorf("failed to connect to database: %w", err)
//		}
//
//		repo := repository.New(db)
//		service := service.New(repo)
//		server := server.New(service, cfg.Port)
//
//		// Return the wire bundle with Runnables and CleanupFunc
//		return ezapp.WireBundle{
//			Runnables: []ezapp.Runnable{server},
//			CleanupFunc: func() error {
//				return db.Close()
//			},
//		}, nil
//	}
//
// The WireFunc is called by the Build function after loading configuration from environment variables.
// If the WireFunc returns an error, Build will attempt to handle it with the provided error handler.
type WireFunc[C any] func(C) (WireBundle, error)
