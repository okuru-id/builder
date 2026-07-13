package seeders

import (
	"encoding/json"

	"gorm.io/datatypes"

	"okuru/app/facades"
	"okuru/app/models"
)

type landingComponentSeed struct {
	Key  string
	Name string
	Tree map[string]any
}

func seedLandingComponents() error {
	seeds := []landingComponentSeed{
		{Key: "hero", Name: "Hero", Tree: wrapRoot(heroSection())},
		{Key: "cta-button", Name: "CTA Button", Tree: wrapRoot(node("button", "CTA Button", map[string]any{"text": "Get started"}, []string{"inline-flex", "items-center", "justify-center", "rounded-md", "bg-blue-600", "px-5", "py-3", "font-medium", "text-white"}))},
		{Key: "feature-grid", Name: "Feature Grid", Tree: wrapRoot(featureGridSection())},
		{Key: "testimonial", Name: "Testimonial", Tree: wrapRoot(testimonialSection())},
		{Key: "faq", Name: "FAQ", Tree: wrapRoot(faqSection())},
		{Key: "pricing-cards", Name: "Pricing Cards", Tree: wrapRoot(pricingSection())},
		{Key: "stats-strip", Name: "Stats Strip", Tree: wrapRoot(statsSection())},
		{Key: "logo-cloud", Name: "Logo Cloud", Tree: wrapRoot(logoCloudSection())},
		{Key: "contact-form", Name: "Contact Form", Tree: wrapRoot(contactFormSection())},
		{Key: "footer-cta", Name: "Footer CTA", Tree: wrapRoot(footerCTASection())},
		{Key: "navbar-simple", Name: "Navbar Simple", Tree: wrapRoot(navbarSection())},
		{Key: "gallery-grid", Name: "Gallery Grid", Tree: wrapRoot(gallerySection())},
	}

	for _, seed := range seeds {
		b, _ := json.Marshal(seed.Tree)
		var comp models.LandingComponent
		if err := facades.Orm().Query().UpdateOrCreate(
			&comp,
			models.LandingComponent{Key: seed.Key},
			models.LandingComponent{Name: seed.Name, Key: seed.Key, IsSystem: true, Tree: datatypes.JSON(b)},
		); err != nil {
			return err
		}
	}
	return nil
}

func wrapRoot(root map[string]any) map[string]any {
	return map[string]any{"root": root}
}

func node(kind, name string, props map[string]any, classes []string, children ...map[string]any) map[string]any {
	if props == nil {
		props = map[string]any{}
	}
	if classes == nil {
		classes = []string{}
	}
	return map[string]any{
		"id":       name,
		"type":     kind,
		"name":     name,
		"props":    props,
		"classes":  classes,
		"children": children,
	}
}

func heroSection() map[string]any {
	return node("section", "Hero", nil, []string{"px-8", "py-20", "bg-neutral-950", "text-white"},
		node("frame", "Hero Content", nil, []string{"mx-auto", "max-w-5xl", "space-y-6", "text-center"},
			node("text", "Eyebrow", map[string]any{"text": "Ship faster"}, []string{"text-sm", "font-medium", "uppercase", "tracking-[0.2em]", "text-neutral-300"}),
			node("heading", "Title", map[string]any{"text": "Build landing pages without the busywork", "level": 1}, []string{"text-4xl", "font-bold", "tracking-tight", "md:text-6xl"}),
			node("text", "Description", map[string]any{"text": "Compose sections, reuse components, and publish polished pages in minutes."}, []string{"mx-auto", "max-w-2xl", "text-base", "text-neutral-300", "md:text-lg"}),
			node("frame", "Actions", nil, []string{"flex", "justify-center", "gap-3", "pt-2", "flex-wrap"},
				node("button", "Primary CTA", map[string]any{"text": "Start free"}, []string{"rounded-md", "bg-white", "px-5", "py-3", "font-medium", "text-neutral-950"}),
				node("button", "Secondary CTA", map[string]any{"text": "Book demo"}, []string{"rounded-md", "border", "border-white/20", "px-5", "py-3", "font-medium", "text-white"}),
			),
		),
	)
}

func featureGridSection() map[string]any {
	return node("section", "Feature Grid", nil, []string{"px-8", "py-16", "bg-white"},
		node("frame", "Container", nil, []string{"mx-auto", "max-w-6xl", "space-y-10"},
			node("heading", "Title", map[string]any{"text": "Everything you need to launch", "level": 2}, []string{"text-3xl", "font-bold", "tracking-tight", "text-center"}),
			node("grid", "Features", nil, []string{"grid", "gap-6", "md:grid-cols-3"},
				featureCard("Fast setup", "Start from reusable building blocks and ship faster."),
				featureCard("Reusable sections", "Save proven sections once, insert them everywhere."),
				featureCard("Publish instantly", "Generate clean output and push updates without fuss."),
			),
		),
	)
}

func testimonialSection() map[string]any {
	return node("section", "Testimonial", nil, []string{"px-8", "py-16", "bg-neutral-50"},
		node("frame", "Quote Wrap", nil, []string{"mx-auto", "max-w-3xl", "space-y-6", "text-center"},
			node("text", "Quote", map[string]any{"text": "“This builder cut our launch time from days to hours.”"}, []string{"text-2xl", "font-medium", "leading-relaxed", "text-neutral-900"}),
			node("text", "Author", map[string]any{"text": "Rina — Product Lead"}, []string{"text-sm", "text-neutral-500"}),
		),
	)
}

func faqSection() map[string]any {
	return node("section", "FAQ", nil, []string{"px-8", "py-16", "bg-white"},
		node("frame", "Container", nil, []string{"mx-auto", "max-w-4xl", "space-y-4"},
			node("heading", "Title", map[string]any{"text": "Frequently asked questions", "level": 2}, []string{"mb-6", "text-3xl", "font-bold"}),
			faqItem("Can I reuse sections?", "Yes. Save any selected node as a reusable component."),
			faqItem("Can I edit copies?", "Duplicate a default component to make an editable custom version."),
			faqItem("Can I publish quickly?", "Yes. Draft, preview, then publish when ready."),
		),
	)
}

func pricingSection() map[string]any {
	return node("section", "Pricing Cards", nil, []string{"px-8", "py-16", "bg-neutral-950", "text-white"},
		node("frame", "Container", nil, []string{"mx-auto", "max-w-6xl", "space-y-10"},
			node("heading", "Title", map[string]any{"text": "Simple pricing", "level": 2}, []string{"text-center", "text-3xl", "font-bold"}),
			node("grid", "Plans", nil, []string{"grid", "gap-6", "md:grid-cols-3"},
				pricingCard("Starter", "$19", "For small launches"),
				pricingCard("Growth", "$49", "For scaling teams"),
				pricingCard("Scale", "$99", "For bigger operations"),
			),
		),
	)
}

func statsSection() map[string]any {
	return node("section", "Stats Strip", nil, []string{"px-8", "py-10", "bg-blue-600", "text-white"},
		node("grid", "Stats", nil, []string{"mx-auto", "grid", "max-w-5xl", "gap-6", "text-center", "md:grid-cols-3"},
			statItem("120+", "teams launched"),
			statItem("3x", "faster iteration"),
			statItem("99.9%", "uptime confidence"),
		),
	)
}

func logoCloudSection() map[string]any {
	return node("section", "Logo Cloud", nil, []string{"px-8", "py-12", "bg-white"},
		node("frame", "Container", nil, []string{"mx-auto", "max-w-5xl", "space-y-6", "text-center"},
			node("text", "Label", map[string]any{"text": "Trusted by fast-moving teams"}, []string{"text-sm", "font-medium", "uppercase", "tracking-[0.2em]", "text-neutral-500"}),
			node("grid", "Logos", nil, []string{"grid", "gap-4", "md:grid-cols-4"},
				logoItem("Acme"),
				logoItem("Nova"),
				logoItem("Pixel"),
				logoItem("Orbit"),
			),
		),
	)
}

func contactFormSection() map[string]any {
	return node("section", "Contact Form", nil, []string{"px-8", "py-16", "bg-neutral-50"},
		node("frame", "Container", nil, []string{"mx-auto", "max-w-3xl", "space-y-6"},
			node("heading", "Title", map[string]any{"text": "Talk to our team", "level": 2}, []string{"text-3xl", "font-bold"}),
			node("form", "Form", map[string]any{"action": "#", "method": "post"}, []string{"space-y-4", "rounded-2xl", "bg-white", "p-6", "shadow-sm"},
				node("input", "Name", map[string]any{"label": "Name", "placeholder": "Your name", "inputType": "text", "required": true}, []string{"w-full"}),
				node("input", "Email", map[string]any{"label": "Email", "placeholder": "you@example.com", "inputType": "email", "required": true}, []string{"w-full"}),
				node("input", "Message", map[string]any{"label": "Message", "placeholder": "Tell us about your project", "inputType": "text", "required": true}, []string{"w-full"}),
				node("button", "Submit", map[string]any{"text": "Send inquiry"}, []string{"rounded-md", "bg-neutral-950", "px-5", "py-3", "font-medium", "text-white"}),
			),
		),
	)
}

func footerCTASection() map[string]any {
	return node("section", "Footer CTA", nil, []string{"px-8", "py-16", "bg-neutral-950", "text-white"},
		node("frame", "Container", nil, []string{"mx-auto", "max-w-4xl", "space-y-4", "text-center"},
			node("heading", "Title", map[string]any{"text": "Ready to launch your next page?", "level": 2}, []string{"text-3xl", "font-bold"}),
			node("text", "Description", map[string]any{"text": "Use the builder, customize sections, and publish when it feels right."}, []string{"text-neutral-300"}),
			node("button", "CTA", map[string]any{"text": "Create page"}, []string{"mx-auto", "inline-flex", "rounded-md", "bg-white", "px-5", "py-3", "font-medium", "text-neutral-950"}),
		),
	)
}

func navbarSection() map[string]any {
	return node("section", "Navbar Simple", nil, []string{"px-8", "py-4", "bg-white", "border-b", "border-neutral-200"},
		node("frame", "Container", nil, []string{"mx-auto", "flex", "max-w-6xl", "items-center", "justify-between", "gap-4"},
			node("text", "Brand", map[string]any{"text": "Brand"}, []string{"text-lg", "font-semibold", "text-neutral-950"}),
			node("frame", "Links", nil, []string{"flex", "items-center", "gap-4", "text-sm", "text-neutral-600"},
				node("link", "Features", map[string]any{"text": "Features", "href": "#features"}, []string{}),
				node("link", "Pricing", map[string]any{"text": "Pricing", "href": "#pricing"}, []string{}),
				node("link", "Contact", map[string]any{"text": "Contact", "href": "#contact"}, []string{}),
			),
		),
	)
}

func gallerySection() map[string]any {
	return node("section", "Gallery Grid", nil, []string{"px-8", "py-16", "bg-white"},
		node("frame", "Container", nil, []string{"mx-auto", "max-w-6xl", "space-y-8"},
			node("heading", "Title", map[string]any{"text": "Show the product in action", "level": 2}, []string{"text-3xl", "font-bold"}),
			node("grid", "Gallery", nil, []string{"grid", "gap-4", "md:grid-cols-3"},
				galleryItem(),
				galleryItem(),
				galleryItem(),
			),
		),
	)
}

func featureCard(title, desc string) map[string]any {
	return node("frame", title, nil, []string{"rounded-2xl", "border", "border-neutral-200", "p-6", "space-y-3", "bg-white"},
		node("heading", "Feature Title", map[string]any{"text": title, "level": 3}, []string{"text-xl", "font-semibold", "text-neutral-950"}),
		node("text", "Feature Description", map[string]any{"text": desc}, []string{"text-sm", "text-neutral-600"}),
	)
}

func faqItem(question, answer string) map[string]any {
	return node("frame", question, nil, []string{"rounded-xl", "border", "border-neutral-200", "p-5", "space-y-2"},
		node("heading", "Question", map[string]any{"text": question, "level": 3}, []string{"text-lg", "font-semibold"}),
		node("text", "Answer", map[string]any{"text": answer}, []string{"text-sm", "text-neutral-600"}),
	)
}

func pricingCard(plan, price, desc string) map[string]any {
	return node("frame", plan, nil, []string{"rounded-2xl", "bg-white", "p-6", "space-y-4", "text-neutral-950"},
		node("heading", "Plan", map[string]any{"text": plan, "level": 3}, []string{"text-xl", "font-semibold"}),
		node("text", "Price", map[string]any{"text": price}, []string{"text-4xl", "font-bold"}),
		node("text", "Description", map[string]any{"text": desc}, []string{"text-sm", "text-neutral-600"}),
		node("button", "Choose", map[string]any{"text": "Choose plan"}, []string{"rounded-md", "bg-neutral-950", "px-4", "py-2", "text-white"}),
	)
}

func statItem(value, label string) map[string]any {
	return node("frame", label, nil, []string{"space-y-1"},
		node("text", "Value", map[string]any{"text": value}, []string{"text-4xl", "font-bold"}),
		node("text", "Label", map[string]any{"text": label}, []string{"text-sm", "text-blue-100"}),
	)
}

func logoItem(name string) map[string]any {
	return node("frame", name, nil, []string{"rounded-xl", "border", "border-neutral-200", "px-6", "py-5", "text-sm", "font-semibold", "text-neutral-500"},
		node("text", "Logo", map[string]any{"text": name}, []string{"text-center"}),
	)
}

func galleryItem() map[string]any {
	return node("image", "Gallery Item", map[string]any{"src": "https://placehold.co/800x600", "alt": "Gallery image"}, []string{"h-56", "w-full", "rounded-2xl", "object-cover"})
}
