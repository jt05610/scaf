package ts

import "text/template"

func PackageJson() *template.Template {
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

func tsconfigJson() *template.Template {
	t, err := template.New("package.json").Parse(`
{
  "compilerOptions": {
    "strict": true,
    "target": "ESNext",
    "module": "ESNext",
    "moduleResolution": "node",
    "allowSyntheticDefaultImports": true,
    "esModuleInterop": true,
    "jsx": "preserve",
    "jsxImportSource": "solid-js",
    "types": ["vite/client", "@testing-library/jest-dom"]
  }
}
`)
	if err != nil {
		panic(err)
	}
	return t
}

//goland:noinspection HttpUrlsUsage
func viteConfigTs() *template.Template {
	t, err := template.New("package.json").Parse(`
/// <reference types="vitest" />
/// <reference types="vite/client" />

import {defineConfig} from 'vite';
import federation from "@originjs/vite-plugin-federation";
import solidPlugin from 'vite-plugin-solid';

export default defineConfig({
    plugins: [
        solidPlugin(),
        federation({
            name: '{{.Name}}',
			filename: 'remoteEntry.js',
            shared: [
                "solid-js",
			],
        })
    ],
    server: {
        port: {{.Port}},
    },
    test: {
        environment: 'jsdom',
        globals: true,
        transformMode: {web: [/\.[jt]sx?$/]},
        setupFiles: ['node_modules/@testing-library/jest-dom/extend-expect.js'],
        // otherwise, solid would be loaded twice:
        deps: {registerNodeLoader: true},
        threads: false,
        isolate: false,
    },
    build: {
        target: 'esnext',
    },
    resolve: {
        conditions: ['development', 'browser'],
    },
});

`)
	if err != nil {
		panic(err)
	}
	return t
}

func indexHtml() *template.Template {
	t, err := template.New("index.html").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>{{.Name}}</title>
  </head>
  <body>
    <noscript>You need to enable JavaScript to run this app.</noscript>
    <div id="root"></div>

    <script src="/src/index.tsx" type="module"></script>
  </body>
</html>
`)
	if err != nil {
		panic(err)
	}
	return t
}

func indexTsx() *template.Template {
	t, err := template.New("index.tsx").Parse(`
import { render } from 'solid-js/web';

import {App} from './app';

const root = document.getElementById('root');

if (import.meta.env.DEV && !(root instanceof HTMLElement)) {
  throw new Error(
    'Root element not found. Did you forget to add it to your index.html? Or maybe the id attribute got mispelled?',
  );
}

render(() => <App />, root!);
`)
	if err != nil {
		panic(err)
	}
	return t
}

func appTsx() *template.Template {
	t, err := template.New("index.tsx").Parse(`
import { Component } from 'solid-js';
import { createStore } from 'solid-js/store';

type {{.Name}}Props = {
	{{- range .Fields}}
	{{.Name}}: {{.Type}}
	{{- end}}
}

export const App : Component = (props: {{.Name}}Props) => {
	return (
		<div>
			Hello, World!
		</div>
)
}
`)
	if err != nil {
		panic(err)
	}
	return t
}
