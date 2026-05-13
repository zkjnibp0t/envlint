// Package templater renders .env template files by substituting
// {{ PLACEHOLDER }} tokens with values from a provided map.
//
// Templates follow a simple double-brace syntax:
//
//	DB_HOST={{ DB_HOST }}
//	DB_PORT={{ DB_PORT }}
//	APP_ENV=production
//
// Usage:
//
//	values := map[string]string{
//		"DB_HOST": "localhost",
//		"DB_PORT": "5432",
//	}
//	result, err := templater.RenderFile("template.env", values)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(result.Rendered)
//	if len(result.Unresolved) > 0 {
//		fmt.Println("Unresolved placeholders:", result.Unresolved)
//	}
package templater
