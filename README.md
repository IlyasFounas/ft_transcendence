*This project has been created as part of the 42 curriculum by eelissal, ifounas, lserodon, spunyapr.*
# ATLAS Intelligence Network

## Description

**ATLAS Intelligence Network** is a web-based multiplayer geolocation game inspired by GeoGuessr, developed as part of the 42 `ft_transcendence` curriculum.

In this game players take on the role of intelligence agents tasked with locating the origin of a signal before it disappears. Using the visual clues provided by a street-level view, players must analyze their surroundings, identify the location, and make their guess before time runs out.

The application combines geolocation gameplay with social and real-time features, including user profiles, friends, chat, online-status tracking, multiplayer lobbies, game statistics, and leaderboards.

The application is built with **React, TypeScript, and Tailwind CSS** on the frontend, and **Go with Gin** on the backend. **PostgreSQL** is used for data persistence, with **Bun ORM** for database access. WebSockets power real-time features such as multiplayer lobbies, chat, and friend-status updates.

## Instructions & Requirements

### Prerequisites

Make sure you have the following installed:

- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- `make`


### Setup & Installation

#### 1. Clone the repository

```bash
git clone <repository-url>
cd ft_transcendence
```

#### 2. Configure Environment Variables

Create your environment file from the example:

```bash
cp .env.exemple .env
```

Open `.env` and fill in the required values, such as database credentials, authentication secrets, and OAuth credentials if applicable.

For production, configure `docker-config/.env.prod` with the required production environment variables.

#### 3. Start the application

**Development mode**

Start the application using the base and development Docker Compose configurations:

```bash
make up
```

**Production mode**

Start the application using the production configuration:

```bash
make prod
```

### Useful Makefile Commands

#### Application Management

| Command | Description |
|---|---|
| `make up` | Start the application in development mode. |
| `make prod` | Update the production IP configuration and start the production containers. |
| `make down` | Stop and remove the development containers. |
| `make stop` | Stop the development containers without removing them. |
| `make restart` | Restart the development containers. |
| `make ps` | Display the status of the development containers. |
| `make re` | Remove containers, volumes, and images, then start the application again. |

#### Logs & Debugging

| Command | Description |
|---|---|
| `make logs-<service>` | Follow logs for a specific service. |
| `make in-<container>` | Open a shell inside a container. |
| `make print_posgres_ip` | Display IP address information from the PostgreSQL container. |

Examples:

```bash
make logs-golang
make logs-react
make in-golang
```

Replace the service or container name with the appropriate name from your Docker Compose configuration.

#### Database Migrations

ATLAS uses Bun ORM migrations to manage changes to the database schema.

| Command | Description |
|---|---|
| `make init_mig` | Initialize the migration system. |
| `make create_mig NAME=<name>` | Create a new migration file. |
| `make exec_mig` | Apply pending migrations. |
| `make rollback_mig` | Roll back the last migration. |
| `make status_mig` | Display migration status. |

Example:

```bash
make create_mig NAME=add_user_fields
make exec_mig
make status_mig
```

#### Cleanup

| Command | Description |
|---|---|
| `make fclean` | Remove development containers, volumes, and images. |
| `make clean-images` | Clean up unused Docker images after running `make clean`. |
| `make clean-caches` | Clean Docker images and Go/npm caches after running `make clean`. |

**Warning:** `make re`, `make fclean`, and the cleanup commands can delete data or cached dependencies. Use them carefully, especially if you need to preserve local database data.


### Access


## Resources

List of classic references related to the topic (documentation, articles, tutorials, etc.).

### References

- [React Documentation (2026)](https://react.dev/learn): official guide to React, covering components, hooks and context, which I used to structure the pages and share the WebSocket state across the app.
- [TypeScript Documentation (2026)](https://www.typescriptlang.org/docs/): official handbook and reference for TypeScript, used to type the data exchanged between the frontend and the backend.
- [Tailwind CSS Documentation (2026)](https://tailwindcss.com/docs/installation/using-vite): reference for the utility classes used to style every page and build responsive layouts.
- [MDN - CSS container queries (2026)](https://developer.mozilla.org/en-US/docs/Web/CSS/Guides/Containment/Container_queries): explanation of container queries, used to scale the podium according to the size of its container.
- [MDN - WebSockets API (2026)](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API): guide to opening, using and closing a WebSocket connection from the browser, at the core of the real-time features.
- [MapillaryJS Documentation (2026)](https://mapillary.github.io/mapillary-js/docs): documentation of the viewer used to display the street-level images of the game.
- [Leaflet Documentation (2026)](https://leafletjs.com/reference.html): reference for the mapping library used for the guess map and the results maps.


### AI Use

// list of each person 


## Team Information

| Member | Role | Responsibilities |
|---|---|---|
| **eelissal** | Product Owner | Feature prioritization, game logic, and backend development |
| **spunyapr** | Project Manager | Project coordination and frontend/backend development |
| **ifounas** | Backend / DevOps Tech Lead | Backend architecture, authentication, infrastructure, and deployment |
| **lserodon** | Frontend Tech Lead | Frontend architecture/development, and reusable components |

All team members contribute to development, testing, and code review.


## Project Management

### Tools

- **GitHub:** Version control, feature branches, pull requests, and code reviews.
- **Trello:** Task planning, assignment, and progress tracking.
- **Discord:** Team communication and coordination.
- **Google Sheets:** Bug tracking and issue follow-up.

### Workflow

Our team follows a collaborative workflow to organize development, monitor progress, and support each other throughout the project.

- **Task Distribution:** Tasks are assigned according to each member's role, responsibilities, and the project's priorities.
- **Weekly Meetings:** The team meets at least once a week to review each member's progress, discuss completed and ongoing tasks, and identify any difficulties or blockers. When someone encounters a problem, other members can provide support or help find a solution.
- **Weekly Task Updates:** At the beginning of each week, team members share their planned tasks in their dedicated Discord channels. These channels are also used to communicate progress, ask questions, and report difficulties throughout the week.
- **Git Workflow:** Development is organized using Git branches. Changes are submitted through pull requests and reviewed before being merged.
- **Testing and Bug Tracking:** Features are tested during development, and identified bugs are documented in Google Sheets to track their resolution.

## Technical Stack

### Frontend

<p align="left">
  <img src="https://cdn.simpleicons.org/react/61DAFB" width="40" alt="React" />
  <img src="https://cdn.simpleicons.org/typescript/3178C6" width="40" alt="TypeScript" />
  <img src="https://cdn.simpleicons.org/tailwindcss/06B6D4" width="40" alt="Tailwind CSS" />
  <img src="https://cdn.simpleicons.org/vite/646CFF" width="40" alt="Vite" />
</p>

| Technology | Purpose | Why choose this? |
|---|---|---|
| **React** | Building the user interface and reusable components | Its component-based architecture makes development and maintenance easier. |
| **TypeScript** | Adding static typing to JavaScript | Helps catch errors early and improves maintainability. |
| **Tailwind CSS** | Styling the application and creating responsive layouts | Enables fast, consistent UI development with utility classes. |
| **Vite** | Frontend development server and build tool | Provides fast startup, hot module replacement, and efficient builds. |

### Backend

<p align="left">
  <img src="https://cdn.simpleicons.org/go/00ADD8" width="40" alt="Go" />
  <img src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png" width="40" alt="Gin" />
  <img src="https://bun.uptrace.dev/hero/logo.png" width="40" alt="Bun ORM" />
</p>

| Technology | Purpose | Why choose this? |
|---|---|---|
| **Go** | Backend logic and API development | Offers strong performance and efficient concurrency. |
| **Gin** | HTTP routing and API handling | Provides lightweight routing and middleware support. |
| **Bun ORM** | Database access and query management | Simplifies database queries within the Go application. |
| **WebSockets** | Real-time communication | Enables instant updates for chat, lobbies, and online status. |

### Infrastructure & DevOps

<p align="left">
  <img src="https://cdn.simpleicons.org/docker/2496ED" width="40" alt="Docker" />
  <img src="https://cdn.simpleicons.org/nginx/009639" width="40" alt="Nginx" />
  <img src="https://cdn.simpleicons.org/github/FFFFFF" width="40" alt="GitHub" />
</p>

| Technology | Purpose | Why choose this? |
|---|---|---|
| **Docker Compose** | Running and coordinating application services | Simplifies multi-service setup and provides consistent environments. |
| **Nginx** | Web server and reverse proxy | Routes incoming requests to the appropriate services. |
| **GitHub** | Repository hosting and collaboration | Supports branches, pull requests, and code reviews. |

### Database

<p align="left">
  <img src="https://cdn.simpleicons.org/postgresql/4169E1" width="40" alt="PostgreSQL" />
</p>

| Technology | Purpose | Why choose this? |
|---|---|---|
| **PostgreSQL** | Storing accounts, friendships, games, and statistics | Provides reliable relational storage, transactions, and data integrity. |

### Design

<p align="left">
  <img src="https://cdn.simpleicons.org/figma/F24E1E" width="40" alt="Figma" />
</p>

| Technology | Purpose | Why choose this? |
|---|---|---|
| **Figma** | Designing interfaces and prototypes | Helps the team plan layouts and maintain a consistent visual design. |

## Database Schema

ATLAS uses **PostgreSQL** for persistent storage and **Bun ORM** to map Go models to database tables. The database supports user management, multiplayer gameplay, geographical guesses, player progression, social interactions, and authentication.

### Tables and Relationships

The tables below are organized by functional priority, starting with the core game flow and user data.

### 1. Users

The `users` table stores player accounts, profile information, authentication-related data, and progression.

| Field | Type | Constraints / Default | Description |
|---|---|---|---|
| `id` | `BIGINT` | Primary key, auto-increment | Unique user identifier. |
| `login` | `VARCHAR` | Unique, nullable | User login. |
| `name` | `VARCHAR` | Unique, `NOT NULL` | Username. |
| `google_id` | `VARCHAR` | Unique, nullable | Google account identifier. |
| `password` | `VARCHAR` | Nullable | Stores the user's hashed password for authentication. |
| `email` | `VARCHAR`  | Unique, `NOT NULL` | User email address. |
| `is_verified` | `BOOLEAN` | `NOT NULL` | Whether the account is verified. |
| `avatar` | `VARCHAR` | Default: `'/avatars/avatar1.png'` | Avatar path. |
| `xp` | `BIGINT` | `NOT NULL`, default: `0` | Player experience points. |
| `rank` | `VARCHAR` | `NOT NULL`, default: `'Junior Agent'` | Player rank. |
| `is_admin` | `BOOLEAN` | `NOT NULL`, default: `false` | Indicates administrator privileges. |
| `created_at` | `TIMESTAMP` | Default: `CURRENT_TIMESTAMP` | Account creation time. |

**Relationships:** A user has one associated statistics record and can have multiple lobby records, according to the model declarations.

### 2. Game

The `game` table stores game configuration and session information.

| Field | Type | Constraints / Default | Description |
|---|---|---|---|
| `id` | `BIGINT` | Primary key, auto-increment | Unique game identifier. |
| `title` | `VARCHAR` | Nullable | Game title. |
| `mode` | `VARCHAR` | Nullable | Game mode. |
| `difficulty` | `VARCHAR` | Nullable | Game difficulty. |
| `timer` | `BIGINT` | `NOT NULL` | Timer configuration. |
| `round_total` | `BIGINT` | `NOT NULL` | Number of rounds in the game. |
| `finished_at` | `TIMESTAMP` | Nullable | Time the game finished. |
| `player_nb` | `BIGINT` | `NOT NULL`, check: `player_nb <> 0` | Configured player count; zero is prohibited. |
| `created_at` | `TIMESTAMP` | Default: `CURRENT_TIMESTAMP` | Game creation time. |
| `location_ids` | `JSONB` | Nullable | Stores location IDs as JSON data. |

**Relationships:** A game has multiple lobby records. The `Round` model also declares a relationship to a game through `game_id`.

### 3. Lobby

The `lobby` table stores player participation in game lobbies.

| Field | Type | Constraints / Default | Description |
|---|---|---|---|
| `id` | `BIGINT` | Primary key, auto-increment | Unique lobby record identifier. |
| `user_id` | `BIGINT` | Nullable in the model | ID of the associated user. |
| `game_id` | `BIGINT` | Nullable in the model | ID of the associated game. |
| `join_code` | `VARCHAR` | Unique, nullable | Code used to join a lobby. |
| `is_creator` | `BOOLEAN` | No explicit default | Indicates whether the user created the lobby. |
| `left_at` | `TIMESTAMP` | Nullable | Time the user left the lobby. |
| `kicked` | `BOOLEAN` | `NOT NULL`, default: `false` | Indicates whether the player was kicked and cannot rejoin the game. |

**Relationships:** The model declares a relationship between `lobby.game_id` and `game.id`. The `User` model also declares a relationship using `user_id`.

### 4. Rounds

The `rounds` table stores individual rounds within a game.

| Field | Type | Constraints / Default | Description |
|---|---|---|---|
| `id` | `BIGINT` | Primary key, auto-increment | Unique round identifier. |
| `round_nb` | `BIGINT` | `NOT NULL` | Round number. |
| `game_id` | `BIGINT` | `NOT NULL` | Associated game ID. |
| `loc_id` | `BIGINT` | `NOT NULL` | ID of the location used for the round. |
| `started_at` | `TIMESTAMP` | `NOT NULL` | Round start time. |

**Relationships:** Each round belongs to a game and a location, according to the model declarations.

### 5. Guesses

The `guesses` table stores player guesses and their results for individual rounds.

| Field | Type | Constraints / Default | Description |
|---|---|---|---|
| `guess_id` | `BIGINT` | Primary key, auto-increment | Unique guess identifier. |
| `latitude` | `DOUBLE PRECISION` | Nullable | Latitude submitted by the player. |
| `longitude` | `DOUBLE PRECISION` | Nullable | Longitude submitted by the player. |
| `country_code` | `VARCHAR(2)` | Nullable | Country code associated with the result. |
| `distance` | `BIGINT` | Nullable | Distance between the guess and target. |
| `country_found` | `BOOLEAN` | Default: `false` | Indicates whether the country was found. |
| `time_left` | `BIGINT` | Default: `0` | Remaining time. |
| `streak` | `BIGINT` | Default: `0` | Current streak value. |
| `base_score` | `BIGINT` | `NOT NULL` | Score before multipliers. |
| `score` | `BIGINT` | `NOT NULL` | Final score after multipliers. |
| `round_nb` | `BIGINT` | `NOT NULL` | Round number. |
| `round_id` | `BIGINT` | `NOT NULL`, composite unique index | Associated round ID. |
| `user_id` | `BIGINT` | `NOT NULL`, composite unique index | User who made the guess. |

**Constraints and relationships:** `round_id` and `user_id` share the composite unique index `idx_guess_user_round`, preventing duplicate guess records for the same user and round. Each guess belongs to a user and a round in the model declarations.

### 6. Locations

The `locations` table stores geographical data used in gameplay.

| Field | Type | Constraints / Default | Description |
|---|---|---|---|
| `location_id` | `BIGINT` | Primary key, auto-increment | Unique location identifier. |
| `mapillary_id` | `BIGINT` | Unique, `NOT NULL` | Mapillary identifier for the location. |
| `latitude` | `DOUBLE PRECISION` | `NOT NULL` | Location latitude. |
| `longitude` | `DOUBLE PRECISION` | `NOT NULL` | Location longitude. |
| `country_code` | `VARCHAR(2)` | `NOT NULL` | Two-character country code. |

### 7. Statistics

The `statistics` table stores player performance metrics.

| Field | Type | Constraints / Default | Description |
|---|---|---|---|
| `id` | `BIGINT` | Primary key, auto-increment | Unique statistics record identifier. |
| `user_id` | `BIGINT` | `NOT NULL`, unique | Associated user ID. |
| `finished_games` | `BIGINT` | `NOT NULL`, default: `0` | Number of finished games. |
| `accuracy` | `DOUBLE PRECISION` | `NOT NULL`, default: `0` | Player accuracy metric. |
| `best_score` | `BIGINT` | `NOT NULL`, default: `0` | Best recorded score. |
| `closest_guess` | `BIGINT` | Nullable | Closest recorded guess value. |

**Relationships:** The unique `user_id` and the model's `has-one` declaration represent one statistics record per user.

### 8. Friends

The `friends` table stores friendship requests between users.

| Field | Type | Constraints / Default | Description |
|---|---|---|---|
| `id` | `BIGINT` | Primary key, auto-increment | Unique friendship record identifier. |
| `requester_id` | `BIGINT` | `NOT NULL` | User who sent the request. |
| `receiver_id` | `BIGINT` | `NOT NULL` | User receiving the request. |
| `status` | `VARCHAR` | `NOT NULL`, default: `'pending'` | Friendship request status. |
| `created_at` | `TIMESTAMP` | Default: `CURRENT_TIMESTAMP` | Record creation time. |

### 9. Chats

The `chats` table stores messages exchanged between users.

| Field | Type | Constraints / Default | Description |
|---|---|---|---|
| `id` | `BIGINT` | Primary key, auto-increment | Unique message identifier. |
| `src_user_id` | `BIGINT` | `NOT NULL` | Sender's user ID. |
| `dst_user_id` | `BIGINT` | `NOT NULL` | Recipient's user ID. |
| `message` | `VARCHAR` | `NOT NULL` | Message content. |
| `created_at` | `TIMESTAMP` | Default: `CURRENT_TIMESTAMP` | Message creation time. |

### 10. Activation Tokens

The `activation_tokens` table stores tokens used for account activation.

| Field | Type | Constraints / Default | Description |
|---|---|---|---|
| `id` | `BIGINT` | Primary key, auto-increment | Unique token record identifier. |
| `user_id` | `BIGINT` | `NOT NULL` | Associated user ID. |
| `token` | `VARCHAR` | `NOT NULL`, unique | Activation token. |
| `expires_at` | `TIMESTAMP` | `NOT NULL` | Token expiration time. |
| `used_at` | `TIMESTAMP` | Nullable | Time the token was used. |

**Relationships:** The model declares a relationship to `users` through `user_id`.

### 11. Refresh Tokens

The `refresh_tokens` table stores hashed tokens used to refresh authentication sessions.

| Field | Type | Constraints / Default | Description |
|---|---|---|---|
| `id` | `BIGINT` | Primary key, auto-increment | Unique token record identifier. |
| `hash_token` | `VARCHAR` | `NOT NULL`, unique | Hashed refresh token. |
| `user_id` | `BIGINT` | `NOT NULL` | Associated user ID. |
| `expires_at` | `TIMESTAMP` | `NOT NULL` | Token expiration time. |
| `used_at` | `TIMESTAMP` | Nullable | Time the token was used. |

**Relationships:** The model declares a relationship to `users` through `user_id`.

### 12. Migration Tables

Bun uses these tables to manage database schema changes.

| Table | Description |
|---|---|
| `bun_migrations` | Records migrations that have been applied. |
| `bun_migration_locks` | Supports locking during migration operations. |

### Relationships Overview

// put database schema diagram

## Features List

| Feature | Description | Contributors |
|---|---|---|
| User Authentication | User registration, login, account verification, and authentication management. | ifounas, spunyapr |
| Google OAuth | Authentication through Google accounts. | ifounas|
| User Profiles | User profile information, avatars, and profile management. | ifounas, lserodon , spunyapr |
| Friends System | Send, receive, and manage friend requests. | spunyapr |
| Friend Status | Display friends' online/offline status and availability in real time. | ifounas, spunyapr |
| Real-Time Chat | Exchange messages with other users. | ifounas, spunyapr |
| Game Lobby | Create and join game lobbies using a join code. | eelissal, ifounas, lserodon |
| Geolocation Gameplay | Identify locations using geographical clues and submit guesses. | eelissal, lserodon |
| Solo Game | Play individually, explore locations, submit guesses, and earn scores without competing against other players. | eelissal, lserodon |
| Multiplayer Games | Play game sessions with other players. | eelissal, ifounas, lserodon |
| Round & Scoring System | Manage game rounds, calculate guess distances, and assign scores. | eelissal, lserodon |
| Player Statistics | Track completed games, accuracy, best score, and closest guess. | eelissal, spunyapr |
| XP & Ranking | Award experience points and track player progression. | eelissal, spunyapr |
| Leaderboard | Compare player scores and rankings. | eelissal, lserodon |
| Admin Permissions | Manage administrative privileges and protected actions. | eelissal, spunyapr |
| Responsive Design | Build application pages and reusable components for a consistent user experience across different screen sizes. | lserodon, spunyapr |
| Design & User Experience | Create design references and maintain a consistent interface. | lserodon, spunyapr |
| Containerization & Deployment | Containerize the application and manage development and production environments using Docker Compose and Nginx. | ifounas |
| Development Environment Setup | Automate common development tasks using Makefile commands, including starting, stopping, and rebuilding services. | ifounas |
| Infrastructure & Service Management | Configure and manage application services, networking, and deployment infrastructure. | ifounas |


## Modules


## Individual Contributions
