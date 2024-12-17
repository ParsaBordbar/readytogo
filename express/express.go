package express

import (
	"log"
	"os"
	"os/exec"
)

// Express initializes a basic Express.js project structure
func Express(projName string) {
	outputDir := "output"

	// Create or cd to output dir
	createOrNavigate(outputDir)

	// Create or cd to project dir
	createOrNavigate(projName)

	// Install NPM packages
	runCommand("npm", "init", "-y")
	runCommand("npm", "install", "express", "nodemon", "joi", "mongoose", "helmet", "morgan", "axios")

	// Create subdir and files
	subDirs := []string{"routes", "controllers", "middlewares", "models", "configs", "services"}
	createDirectories(subDirs)

	fileNames := []string{"index.js", "app.js"}
	createFiles(fileNames)

    serverJSCode := `
const app = require('./app');
const config = require('./configs/config')
const mongoose = require('mongoose');

mongoose.connect(config.mongoose.url)
    .then(() => console.log("Connected to Mongodb :)"))
    .catch((err) => console.error("Could not connect to Mongodb.", err));

const port = config.port;
app.listen(port, () => console.log("Listening on Port"));
`
writeFile("index.js", serverJSCode)

    appJsCode := `
const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const morgan = require('morgan');
const routes = require('./routes/routes')

const app = express();

app.use(cors());
app.options('*', cors());
app.use(morgan('combined'));
app.use(helmet());

app.use('/api/v1', routes);

module.exports = app;
    `
    writeFile("app.js", appJsCode)

    routesJsCode := `
const express = require('express');
const router = express.Router();

const routes = [
    {
        path: '/path',
        route: route
    }
];

routes.forEach((route) => {
    router.use(route.path, route.route);
});

module.exports = router;
    `
	// Create `routes.js` in the routes directory if not exist
	if err := os.Chdir("routes"); err != nil {
		log.Fatalf("Failed to change directory to routes: %v", err)
	}
	if _, err := os.Stat("routes.js"); os.IsNotExist(err) {
		writeFile("routes.js", routesJsCode)
	}

	log.Println("Express.js project setup completed successfully.")
}

func createOrNavigate(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}
	if err := os.Chdir(dir); err != nil {
		log.Fatalf("Failed to change directory to %s: %v", dir, err)
	}
}

func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to execute command %s: %v", name, err)
	}
}

func createDirectories(dirs []string) {
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.Mkdir(dir, 0755); err != nil {
				log.Printf("Failed to create directory %s: %v", dir, err)
			}
		}
	}
}

func createFiles(files []string) {
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			if _, err := os.Create(file); err != nil {
				log.Printf("Failed to create file: %v", err)
			}
		}
	}
}

func writeFile(fileName string, content string) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open file %s: %v", fileName, err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		log.Fatalf("Failed to write to file %s: %v", fileName, err)
	}
}
