package newspaper

// Headline represents a news headline entry for the infinite scroll feed.
type Headline struct {
	ID                               int
	Category, Title, Summary, Byline string
}

// Headlines is the full set of headline entries for the newspaper guide.
var Headlines = []Headline{
	{0, "Design", "The Grid Is Dead, Long Live the Grid", "Modern layout systems have made the rigid grid obsolete — or have they? A look at the evolution of page structure.", "By Jane Chen · 8 min read"},
	{1, "Typography", "Why Your Font Choice Is Wrong", "A provocative look at the assumptions designers make about typeface selection and readability.", "By Marcus Webb · 5 min read"},
	{2, "Color Theory", "The Case Against Color", "When restraint becomes the most powerful tool in a designer's arsenal.", "By Sarah Kim · 6 min read"},
	{3, "CSS", "Container Queries Changed Everything", "How the newest CSS specification is reshaping component-driven design.", "By Dev Patel · 7 min read"},
	{4, "Editorial", "Print Is Not Dead, It Evolved", "The newspaper aesthetic finds new life in digital interfaces.", "By The Editors · 4 min read"},
	{5, "Architecture", "Whitespace Is Not Empty Space", "Understanding the active role of negative space in visual hierarchy.", "By Yuki Tanaka · 5 min read"},
	{6, "Web Standards", "The Semantic Web We Were Promised", "Two decades later, are we any closer to the original vision?", "By Alex Rivera · 9 min read"},
	{7, "Design Systems", "One Component to Rule Them All", "The pursuit of the perfect reusable component — and why it's a trap.", "By Priya Sharma · 6 min read"},
	{8, "Typography", "The Golden Ratio Is Overrated", "Mathematical beauty does not always equal visual beauty.", "By Marcus Webb · 4 min read"},
	{9, "Accessibility", "Designing for Everyone Means Designing for No One", "A counterpoint to universal design — and why specificity matters.", "By Jordan Lee · 7 min read"},
	{10, "CSS", "Flexbox vs Grid: The Final Answer", "Spoiler: the answer is both. But knowing when to use which is the real skill.", "By Dev Patel · 5 min read"},
	{11, "Editorial", "The Attention Economy Broke Design", "How metrics-driven design is undermining craft.", "By The Editors · 3 min read"},
	{12, "Color Theory", "Red Means Stop (Except When It Doesn't)", "Cultural context and the unreliability of color as communication.", "By Sarah Kim · 6 min read"},
	{13, "Architecture", "Every Layout Is a Compromise", "The tensions between content, aesthetics, and engineering.", "By Yuki Tanaka · 8 min read"},
	{14, "Web Standards", "HTML Is a Programming Language", "A deliberately provocative position, rigorously defended.", "By Alex Rivera · 5 min read"},
}

// ArticleData represents a full article with body text.
type ArticleData struct {
	Category, Title, Byline, Body string
}

// Articles maps article IDs to their full content.
var Articles = map[string]ArticleData{
	"0": {"Design", "The Grid Is Dead, Long Live the Grid", "By Jane Chen · March 7, 2026", "The grid has been the backbone of graphic design since the Bauhaus movement. For nearly a century, designers have relied on invisible lines to create order from chaos. But as digital interfaces have grown more fluid and responsive, the rigid grid has begun to feel like a constraint rather than a tool. Modern CSS layout systems — Flexbox, Grid, and now container queries — have given designers unprecedented freedom. Yet paradoxically, this freedom has led many back to the grid, not as a cage, but as a starting point. The best modern layouts use the grid as a foundation, then deliberately break it to create visual tension and hierarchy. The grid is dead. Long live the grid."},
	"1": {"Typography", "Why Your Font Choice Is Wrong", "By Marcus Webb · March 6, 2026", "Every designer has a favorite typeface. For some it is Helvetica, that Swiss army knife of type. For others, it is something more expressive — a Didot, perhaps, or a carefully crafted variable font. But here is the uncomfortable truth: your font choice probably matters less than you think. Research consistently shows that readers adapt to virtually any well-set typeface within seconds. What matters far more is the typographic system — the relationships between sizes, weights, and spacing. A mediocre font set beautifully will always outperform a beautiful font set poorly. Stop agonizing over the typeface. Start obsessing over the system."},
	"2": {"Color Theory", "The Case Against Color", "By Sarah Kim · March 5, 2026", "In a world of vibrant gradients and bold color palettes, there is something radical about restraint. The most powerful designs often use color sparingly — a single accent against a field of neutrals. This newspaper-inspired aesthetic proves the point: with just cream, black, and a touch of red, we can create hierarchy, emphasis, and emotional resonance. Color is not decoration. It is signal. And when everything is colorful, nothing stands out. The next time you reach for a rainbow palette, ask yourself: what if I used just one color instead?"},
	"3": {"CSS", "Container Queries Changed Everything", "By Dev Patel · March 4, 2026", "For years, responsive design meant media queries — asking the viewport how wide it was, then making decisions based on that answer. But components do not live in viewports. They live in containers. A card might appear in a sidebar, a main column, or a modal, each with different available widths. Container queries finally let us ask the right question: how much space does my parent give me? This changes everything about how we think about component design. No more breakpoint gymnastics. No more wrapper divs to simulate container awareness. Just components that know their context and respond accordingly."},
	"4": {"Editorial", "Print Is Not Dead, It Evolved", "By The Editors · March 3, 2026", "Every few years, someone declares print dead. And every few years, print proves them wrong — not by staying the same, but by evolving. The newspaper aesthetic you see on this page is not nostalgia. It is a recognition that centuries of typographic refinement produced principles that transcend medium. Column layouts, drop caps, pull quotes, careful leading — these are not print artifacts. They are solutions to the universal problem of making text readable and engaging. The web did not kill print. It gave print new life."},
}
