# Admin Panel for LabraGo

## Prerequisites
To run this project, ensure you have the following installed on your system:

- [Node.js](https://nodejs.org/) (version 18 or later is recommended)
- [Yarn](https://yarnpkg.com/) (version 1.22 or later)

## Installation
Follow these steps to set up the project on your local machine:

1. Clone the repository:
   ```bash
   git clone https://github.com/GoLabra/labrago.git
   cd labrago
   ```

2. Navigate to the admin project directory:
   ```bash
   cd src/admin
   ```   

3. Install dependencies using Yarn:
   ```bash
   yarn install
   ```

## Running in Development Mode
To run the app in development mode:

1. Start the development server:
   ```bash
   yarn dev
   ```

2. Open your browser and navigate to `http://localhost:3000` to view the app.

## Building the Project
To build the app for production:

1. Run the build command:
   ```bash
   yarn build
   ```

2. The production-ready files will be located in the `.next` directory.

## Starting the Production Server
To start the app in production mode:

1. Ensure the app is built (see the "Building the Project" section).

2. Run the production server:
   ```bash
   yarn start
   ```

3. Open your browser and navigate to the production URL (e.g., `http://localhost:3000`).

## Environment Variables
This project requires environment variables to be set up for proper functionality. Create a `.env.local` file in the root directory and define your variables there. Example:

```
NEXT_PUBLIC_BRAND_PRODUCT_NAME = "Labra·GO"
NEXT_PUBLIC_BRAND_COLOR = "blue"

NEXT_PUBLIC_GRAPHQL_API_URL = "http://your-domain.com:4000"
NEXT_PUBLIC_GRAPHQL_QUERY_API_URL = "http://your-domain.com:4000/query"
NEXT_PUBLIC_GRAPHQL_QUERY_SUBSCRIPTION_URL = "ws://your-domain.com:4000/query"
NEXT_PUBLIC_GRAPHQL_QUERY_PLAYGROUND_URL = "http://your-domain.com:4000/playground"
NEXT_PUBLIC_GRAPHQL_ENTITY_API_URL = "http://your-domain.com:4000/entity"
NEXT_PUBLIC_GRAPHQL_ENTITY_SUBSCRIPTION_URL = "ws://your-domain.com:4000/entity"
NEXT_PUBLIC_GRAPHQL_ENTITY_PLAYGROUND_URL = "http://your-domain.com:4000/eplayground"

NEXT_PUBLIC_CENTRIFUGO_URL = "ws://your-domain.com:8000/connection/websocket"
```

Refer to the `.env` file for a complete list of required variables.

Note: Rebuild the application after making changes to environment variables.


## Contributing
We welcome contributions to this project! To contribute:

1. Fork the repository and clone it to your local machine.
2. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. Commit your changes with descriptive messages:
   ```bash
   git commit -m "Add your feature description"
   ```
4. Push your branch:
   ```bash
   git push origin feature/your-feature-name
   ```
5. Open a pull request on GitHub.
