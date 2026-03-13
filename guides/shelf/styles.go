package shelf

// Book represents a book entry in the style guide demo.
type Book struct {
	ID     string
	Title  string
	Author string
	ISBN   string
	Pages  int
	Genre  string
	Year   int
	Rating float64 // 0-5, half-star precision
	Status string  // "reading", "to-read", "finished"
	Tags   []Tag
	Notes  string
	Review string
}

// Tag represents a triple-encoded tag (icon + color + text).
type Tag struct {
	Icon     string // emoji icon
	Label    string
	Category string // "list", "gift", "club", "rating"
}

// BookList represents a named collection of books.
type BookList struct {
	Name  string
	Icon  string
	Books []Book
}

// tagClass returns the CSS modifier class for a tag category.
func tagClass(category string) string {
	switch category {
	case "list":
		return "shf-tag-list"
	case "gift":
		return "shf-tag-gift"
	case "club":
		return "shf-tag-club"
	case "rating":
		return "shf-tag-rating"
	default:
		return ""
	}
}

// statusLabel returns a human-readable label for a book status.
func statusLabel(status string) string {
	switch status {
	case "reading":
		return "Reading"
	case "to-read":
		return "To Read"
	case "finished":
		return "Finished"
	default:
		return status
	}
}

// statusIcon returns an icon for a book status.
func statusIcon(status string) string {
	switch status {
	case "reading":
		return "📖"
	case "to-read":
		return "📋"
	case "finished":
		return "✅"
	default:
		return "📄"
	}
}

// coverColor returns a warm color for a book cover placeholder based on the book ID.
func coverColor(id string) string {
	colors := []string{
		"#8b7355", "#6b7b5e", "#7b6b8a", "#8a7b6b",
		"#5b7a7a", "#7a5b5b", "#6b8b6b", "#8b6b7b",
		"#7a6b5b", "#5b6b8a",
	}
	sum := 0
	for _, c := range id {
		sum += int(c)
	}
	return colors[sum%len(colors)]
}

// AllBooks is the complete set of demo books.
var AllBooks = []Book{
	{ID: "dune", Title: "Dune", Author: "Frank Herbert", ISBN: "978-0-441-17271-9", Pages: 688, Genre: "Science Fiction", Year: 1965, Rating: 5, Status: "finished", Tags: []Tag{{Icon: "⭐", Label: "Favorite", Category: "rating"}}},
	{ID: "project-hail-mary", Title: "Project Hail Mary", Author: "Andy Weir", ISBN: "978-0-593-13520-4", Pages: 496, Genre: "Science Fiction", Year: 2021, Rating: 4.5, Status: "reading", Tags: []Tag{{Icon: "📖", Label: "Book Club - March", Category: "club"}}},
	{ID: "pachinko", Title: "Pachinko", Author: "Min Jin Lee", ISBN: "978-1-455-56393-7", Pages: 490, Genre: "Historical Fiction", Year: 2017, Rating: 4, Status: "reading"},
	{ID: "sapiens", Title: "Sapiens", Author: "Yuval Noah Harari", ISBN: "978-0-062-31609-7", Pages: 464, Genre: "Non-Fiction", Year: 2011, Rating: 0, Status: "to-read"},
	{ID: "piranesi", Title: "Piranesi", Author: "Susanna Clarke", ISBN: "978-1-635-57563-2", Pages: 272, Genre: "Fantasy", Year: 2020, Rating: 0, Status: "to-read", Tags: []Tag{{Icon: "🎁", Label: "Gift: Mom", Category: "gift"}}},
	{ID: "klara", Title: "Klara and the Sun", Author: "Kazuo Ishiguro", ISBN: "978-0-593-31817-1", Pages: 320, Genre: "Literary Fiction", Year: 2021, Rating: 0, Status: "to-read"},
	{ID: "circe", Title: "Circe", Author: "Madeline Miller", ISBN: "978-0-316-55634-7", Pages: 400, Genre: "Fantasy", Year: 2018, Rating: 0, Status: "to-read", Tags: []Tag{{Icon: "🎁", Label: "Gift: Sarah", Category: "gift"}}},
	{ID: "midnight-library", Title: "The Midnight Library", Author: "Matt Haig", ISBN: "978-0-525-55947-4", Pages: 288, Genre: "Fiction", Year: 2020, Rating: 3.5, Status: "finished", Tags: []Tag{{Icon: "🎁", Label: "Gift: Mom", Category: "gift"}}},
	{ID: "anxious-people", Title: "Anxious People", Author: "Fredrik Backman", ISBN: "978-1-501-16017-8", Pages: 352, Genre: "Fiction", Year: 2019, Rating: 4, Status: "finished", Tags: []Tag{{Icon: "🎁", Label: "Gift: Mom", Category: "gift"}}},
	{ID: "song-of-achilles", Title: "The Song of Achilles", Author: "Madeline Miller", ISBN: "978-0-062-06062-3", Pages: 378, Genre: "Fantasy", Year: 2011, Rating: 4.5, Status: "finished", Tags: []Tag{{Icon: "🎁", Label: "Gift: Sarah", Category: "gift"}, {Icon: "⭐", Label: "Favorite", Category: "rating"}}},
}

// DemoLists is the set of book lists for the dashboard demo.
var DemoLists = []BookList{
	{
		Name:  "Currently Reading",
		Icon:  "📖",
		Books: []Book{AllBooks[1], AllBooks[2]}, // Project Hail Mary, Pachinko
	},
	{
		Name:  "To Be Read",
		Icon:  "📋",
		Books: []Book{AllBooks[3], AllBooks[4], AllBooks[5], AllBooks[6]}, // Sapiens, Piranesi, Klara, Circe
	},
	{
		Name:  "Gift Ideas: Mom",
		Icon:  "🎁",
		Books: []Book{AllBooks[4], AllBooks[7], AllBooks[8]}, // Piranesi, Midnight Library, Anxious People
	},
	{
		Name:  "Book Club - March",
		Icon:  "📖",
		Books: []Book{AllBooks[1]}, // Project Hail Mary
	},
}

// BookByID finds a book by ID.
func BookByID(id string) (Book, bool) {
	for _, b := range AllBooks {
		if b.ID == id {
			return b, true
		}
	}
	return Book{}, false
}

// guideStyles returns the guide-specific CSS classes.
func guideStyles() string {
	return `
/* [custom] - warm panel with subtle shadow */
.shf-panel {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-card);
}
/* [custom] - collapsible list header */
.shf-list-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    cursor: pointer;
    font-family: var(--font-display);
    font-size: var(--font-size-heading);
    font-weight: 600;
    color: var(--color-text);
    border-bottom: 1px solid var(--color-border);
    background: var(--color-surface-2);
    border-radius: var(--radius-md) var(--radius-md) 0 0;
    user-select: none;
    transition: background 0.1s;
}
.shf-list-header:hover {
    background: var(--color-border);
}
/* [custom] - chevron rotation for collapse/expand */
.shf-chevron {
    display: inline-flex;
    font-size: 0.75rem;
    color: var(--color-text-muted);
    transition: transform 0.15s;
    flex-shrink: 0;
}
.shf-chevron-open {
    transform: rotate(90deg);
}
/* [custom] - count badge on list headers */
.shf-badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 1.25rem;
    height: 1.25rem;
    font-family: var(--font-body);
    font-size: 0.6875rem;
    font-weight: 600;
    color: var(--color-text-muted);
    background: var(--color-bg);
    border-radius: 10px;
    padding: 0 0.375rem;
}
/* [custom] - compact book row */
.shf-book-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.625rem 1rem;
    border-bottom: 1px solid var(--color-border);
    cursor: pointer;
    transition: background 0.1s;
}
.shf-book-row:last-child {
    border-bottom: none;
}
.shf-book-row:hover {
    background: var(--color-surface-2);
}
/* [custom] - small book cover placeholder */
.shf-cover {
    width: 40px;
    height: 60px;
    border-radius: 3px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.625rem;
    color: rgba(255,255,255,0.7);
    font-family: var(--font-display);
}
/* [custom] - large cover for detail view */
.shf-cover-lg {
    width: 120px;
    height: 180px;
    border-radius: var(--radius-sm);
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.875rem;
    color: rgba(255,255,255,0.7);
    font-family: var(--font-display);
}
/* [custom] - triple-encoded tag pill (icon + color + text) */
.shf-tag {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
    font-family: var(--font-body);
    font-size: 0.6875rem;
    font-weight: 500;
    padding: 0.125rem 0.5rem;
    border-radius: 10px;
    white-space: nowrap;
}
.shf-tag-icon {
    font-size: 0.75rem;
    line-height: 1;
}
.shf-tag-list { color: var(--color-primary); background: color-mix(in srgb, var(--color-primary) 10%, transparent); }
.shf-tag-gift { color: var(--color-accent); background: color-mix(in srgb, var(--color-accent) 10%, transparent); }
.shf-tag-club { color: var(--color-info); background: color-mix(in srgb, var(--color-info) 10%, transparent); }
.shf-tag-rating { color: var(--color-accent-2); background: color-mix(in srgb, var(--color-accent-2) 10%, transparent); }
/* [custom] - star rating display */
.shf-stars {
    display: inline-flex;
    gap: 1px;
    color: var(--color-accent);
    font-size: 0.875rem;
    line-height: 1;
}
.shf-star-empty {
    color: var(--color-border);
}
/* [custom] - warm-styled button */
.shf-btn {
    font-family: var(--font-body);
    font-size: var(--font-size-caption);
    font-weight: 600;
    color: var(--color-primary);
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    padding: 0.4rem 1rem;
    cursor: pointer;
    transition: background 0.1s, color 0.1s, border-color 0.1s;
}
.shf-btn:hover {
    background: var(--color-primary);
    color: #ffffff;
    border-color: var(--color-primary);
}
.shf-btn-primary {
    background: var(--color-primary);
    color: #ffffff;
    border-color: var(--color-primary);
}
.shf-btn-primary:hover {
    background: color-mix(in srgb, var(--color-primary) 80%, black);
    border-color: color-mix(in srgb, var(--color-primary) 80%, black);
}
.shf-btn-danger {
    color: var(--color-danger);
    border-color: var(--color-danger);
}
.shf-btn-danger:hover {
    background: var(--color-danger);
    color: #ffffff;
}
/* [custom] - scan action button */
.shf-btn-scan {
    background: var(--color-primary);
    color: #ffffff;
    border-color: var(--color-primary);
    padding: 0.5rem 1.25rem;
    font-size: var(--font-size-body);
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
}
.shf-btn-scan:hover {
    background: color-mix(in srgb, var(--color-primary) 80%, black);
    border-color: color-mix(in srgb, var(--color-primary) 80%, black);
}
/* [custom] - warm-styled input */
.shf-input {
    font-family: var(--font-body);
    font-size: var(--font-size-body);
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    color: var(--color-text);
    padding: 0.4rem 0.75rem;
    width: 100%;
    transition: border-color 0.15s;
}
.shf-input:focus {
    outline: none;
    border-color: var(--color-primary);
}
.shf-input::placeholder {
    color: var(--color-text-muted);
}
/* [custom] - search bar */
.shf-search {
    font-family: var(--font-body);
    font-size: var(--font-size-body);
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    color: var(--color-text);
    padding: 0.5rem 1rem 0.5rem 2.25rem;
    width: 100%;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%237a6e60' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Ccircle cx='11' cy='11' r='8'/%3E%3Cline x1='21' y1='21' x2='16.65' y2='16.65'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: 0.75rem center;
}
.shf-search:focus {
    outline: none;
    border-color: var(--color-primary);
}
.shf-search::placeholder {
    color: var(--color-text-muted);
}
/* [custom] - book detail panel */
.shf-detail {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-card);
    padding: 1.5rem;
}
/* [custom] - review card */
.shf-review {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-md);
    padding: 1rem;
}
/* [custom] - warm divider */
.shf-divider {
    border-top: 1px solid var(--color-border);
}
/* [custom] - section accent rule */
.shf-section-rule {
    border-top: 2px solid var(--color-primary);
    padding-top: 1.5rem;
    margin-top: 2rem;
}
/* [custom] - avatar placeholder */
.shf-avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    background: var(--color-surface-2);
    border: 1px solid var(--color-border);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.875rem;
    flex-shrink: 0;
}
/* [custom] - metadata label */
.shf-meta-label {
    font-family: var(--font-body);
    font-size: var(--font-size-caption);
    font-weight: 500;
    color: var(--color-text-muted);
    text-transform: uppercase;
    letter-spacing: 0.03em;
}
/* [custom] - reading status indicator */
.shf-status {
    display: inline-flex;
    align-items: center;
    gap: 0.375rem;
    font-family: var(--font-body);
    font-size: 0.6875rem;
    font-weight: 500;
    padding: 0.2rem 0.625rem;
    border-radius: 10px;
}
.shf-status-reading { color: var(--color-info); background: color-mix(in srgb, var(--color-info) 10%, transparent); }
.shf-status-to-read { color: var(--color-text-muted); background: var(--color-surface-2); }
.shf-status-finished { color: var(--color-success); background: color-mix(in srgb, var(--color-success) 10%, transparent); }
/* [custom] - select input */
.shf-select {
    font-family: var(--font-body);
    font-size: var(--font-size-body);
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: var(--radius-sm);
    color: var(--color-text);
    padding: 0.4rem 0.75rem;
    width: 100%;
}
.shf-select:focus {
    outline: none;
    border-color: var(--color-primary);
}
/* [custom] - checkbox and radio custom styling */
.shf-checkbox {
    accent-color: var(--color-primary);
    width: 1rem;
    height: 1rem;
}
`
}
