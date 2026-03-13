package cassette

type spacingStep struct {
	Px      int
	Label   string
	TwClass string
}

var spacingSteps = []spacingStep{
	{4, "4px / 0.25rem", "p-1"},
	{8, "8px / 0.5rem", "p-2"},
	{12, "12px / 0.75rem", "p-3"},
	{16, "16px / 1rem", "p-4"},
	{20, "20px / 1.25rem", "p-5"},
	{24, "24px / 1.5rem", "p-6"},
	{32, "32px / 2rem", "p-8"},
	{40, "40px / 2.5rem", "p-10"},
	{48, "48px / 3rem", "p-12"},
	{64, "64px / 4rem", "p-16"},
	{80, "80px / 5rem", "p-20"},
	{96, "96px / 6rem", "p-24"},
}
