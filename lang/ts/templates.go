package ts

import "html/template"

func packageJson() *template.Template {
	t, err := template.New("package.json").Parse(`
{
  "name": "{{.Name}}",
  "version": "0.1.0",
  "description": "{{.Desc}}",
  "type": "module",
  "scripts": {
    "start": "vite",
    "dev": "vite",
    "build": "vite build",
    "serve": "vite preview",
    "test": "vitest"
  },
  "license": "MIT",
  "devDependencies": {
    "@originjs/vite-plugin-federation": "^1.2.3",
    "@solidjs/testing-library": "^0.6.0",
    "@testing-library/jest-dom": "^5.16.5",
    "@types/testing-library__jest-dom": "^5.14.5",
    "jsdom": "^21.1.0",
    "typescript": "^4.9.5",
    "vite": "^4.1.1",
    "vite-plugin-solid": "^2.5.0",
    "vitest": "^0.28.4"
  },
  "dependencies": {
    "solid-js": "^1.6.10"
  }
}
`)
	if err != nil {
		panic(err)
	}
	return t
}
