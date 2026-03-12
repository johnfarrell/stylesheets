package tracker

// Item represents a trackable skill, project, certification, or challenge.
type Item struct {
	ID           string
	Category     string // "skill", "project", "certification", "challenge"
	Name         string
	Status       string // "complete", "progress", "locked"
	Level        int    // current level or 0
	Target       int    // target level or 0
	Description  string
	Requirements []string
	Unlocks      []string
}

// Items is the full set of tracker entries displayed in the Mission Control guide.
var Items = []Item{
	// Skills (8)
	{ID: "golang", Category: "skill", Name: "Go", Status: "progress", Level: 75, Target: 99, Description: "Statically typed compiled language for backend services.", Requirements: nil, Unlocks: []string{"Backend APIs", "CLI tools", "Microservices"}},
	{ID: "javascript", Category: "skill", Name: "JavaScript", Status: "progress", Level: 82, Target: 99, Description: "Dynamic language for web frontends and Node.js backends.", Requirements: nil, Unlocks: []string{"React/Vue apps", "Node.js servers"}},
	{ID: "python", Category: "skill", Name: "Python", Status: "complete", Level: 70, Target: 70, Description: "General-purpose language for scripting and data science.", Requirements: nil, Unlocks: []string{"Automation scripts", "Data pipelines"}},
	{ID: "rust", Category: "skill", Name: "Rust", Status: "progress", Level: 80, Target: 99, Description: "Systems language with memory safety guarantees.", Requirements: nil, Unlocks: []string{"Systems programming", "WebAssembly"}},
	{ID: "sql", Category: "skill", Name: "SQL", Status: "progress", Level: 52, Target: 77, Description: "Query language for relational databases.", Requirements: nil, Unlocks: []string{"Database design", "Query optimization", "Migrations"}},
	{ID: "docker", Category: "skill", Name: "Docker", Status: "progress", Level: 85, Target: 99, Description: "Container platform for packaging applications.", Requirements: nil, Unlocks: []string{"Container orchestration", "CI/CD pipelines"}},
	{ID: "git", Category: "skill", Name: "Git", Status: "progress", Level: 72, Target: 85, Description: "Distributed version control system.", Requirements: nil, Unlocks: []string{"Branch strategies", "Rebasing workflows"}},
	{ID: "typescript", Category: "skill", Name: "TypeScript", Status: "complete", Level: 70, Target: 70, Description: "Typed superset of JavaScript.", Requirements: nil, Unlocks: []string{"Type-safe frontends", "Shared API types"}},
	// Projects (6)
	{ID: "todo-cli", Category: "project", Name: "CLI Todo App", Status: "complete", Level: 0, Target: 0, Description: "Build a command-line task manager with file persistence.", Requirements: nil, Unlocks: []string{"CLI patterns", "File I/O experience"}},
	{ID: "rest-api", Category: "project", Name: "REST API", Status: "complete", Level: 0, Target: 0, Description: "Design and implement a RESTful API with authentication.", Requirements: []string{"Go proficiency", "SQL basics", "HTTP fundamentals"}, Unlocks: []string{"API design patterns", "Auth flows"}},
	{ID: "chat-app", Category: "project", Name: "Real-time Chat", Status: "progress", Level: 0, Target: 0, Description: "WebSocket-based chat application with rooms.", Requirements: []string{"JavaScript proficiency", "REST API project", "Basic networking"}, Unlocks: []string{"Real-time protocols", "Event-driven design"}},
	{ID: "blog-engine", Category: "project", Name: "Blog Engine", Status: "locked", Level: 0, Target: 0, Description: "Full-stack blog with SSR, markdown, and comments.", Requirements: []string{"Go proficiency", "Database design", "REST API project", "Docker basics"}, Unlocks: []string{"Full-stack patterns", "SSR experience"}},
	{ID: "search-engine", Category: "project", Name: "Search Engine", Status: "locked", Level: 0, Target: 0, Description: "Build a basic search engine with indexing and ranking.", Requirements: []string{"Go or Rust proficiency", "Data structures", "File I/O", "CLI Todo App project", "REST API project"}, Unlocks: []string{"Information retrieval", "Indexing algorithms"}},
	{ID: "compiler", Category: "project", Name: "Toy Compiler", Status: "locked", Level: 0, Target: 0, Description: "Write a compiler for a small programming language.", Requirements: []string{"Rust proficiency", "Data structures", "Parsing theory"}, Unlocks: []string{"Language design", "Code generation"}},
	// Certifications (3)
	{ID: "aws-ccp", Category: "certification", Name: "AWS Cloud Practitioner", Status: "complete", Level: 0, Target: 0, Description: "Foundational AWS cloud certification.", Requirements: []string{"Cloud computing basics"}, Unlocks: []string{"AWS fundamentals", "Cloud vocabulary"}},
	{ID: "aws-saa", Category: "certification", Name: "AWS Solutions Architect", Status: "progress", Level: 0, Target: 0, Description: "Associate-level AWS architecture certification.", Requirements: []string{"AWS Cloud Practitioner", "Networking basics", "Security fundamentals"}, Unlocks: []string{"Architecture patterns", "AWS service mastery"}},
	{ID: "k8s-cka", Category: "certification", Name: "CKA (Kubernetes)", Status: "locked", Level: 0, Target: 0, Description: "Certified Kubernetes Administrator exam.", Requirements: []string{"Docker proficiency", "Linux administration", "Networking"}, Unlocks: []string{"K8s cluster management", "Container orchestration"}},
	// Challenges (3)
	{ID: "advent-of-code", Category: "challenge", Name: "Advent of Code", Status: "complete", Level: 0, Target: 0, Description: "Annual 25-day coding challenge event.", Requirements: []string{"Any programming language"}, Unlocks: []string{"Algorithm practice", "Problem-solving skills"}},
	{ID: "leetcode-75", Category: "challenge", Name: "LeetCode 75", Status: "progress", Level: 0, Target: 0, Description: "Curated list of 75 essential algorithm problems.", Requirements: []string{"Data structures knowledge", "Algorithm basics"}, Unlocks: []string{"Interview readiness", "Pattern recognition"}},
	{ID: "system-design", Category: "challenge", Name: "System Design", Status: "locked", Level: 0, Target: 0, Description: "Design scalable distributed systems.", Requirements: []string{"Networking", "Database design", "Docker proficiency"}, Unlocks: []string{"Architecture skills", "Senior-level interviews"}},
}
