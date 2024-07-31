import franken from 'franken-ui/shadcn-ui/preset-quick';

/** @type {import('tailwindcss').Config} */
export default {
  presets: [
    franken({
        palette: {
			":root": {
          "--background": "0 0% 100%",
          "--foreground": "0 0% 3.9%",
          "--card": "0 0% 100%",
          "--card-foreground": "0 0% 3.9%",
          "--popover": "0 0% 100%",
          "--popover-foreground": "0 0% 3.9%",
          "--primary": "240 5.9% 10%",
          "--primary-foreground": "0 0% 98%",
          "--secondary": "207 90% 54%", // Changed to a blue color (approximate)
          "--secondary-foreground": "0 0% 9%",
          "--muted": "0 0% 96.1%",
          "--muted-foreground": "0 0% 45.1%",
          "--accent": "0 0% 98%",
          "--accent-foreground": "0 0% 9%",
          "--destructive": "0 84.2% 60.2%",
          "--destructive-foreground": "0 0% 98%",
          "--border": "0 0% 89.8%",
          "--input": "0 0% 89.8%",
          "--ring": "0 0% 3.9%",
        },
        ".dark": {
          "--background": "0 0% 3.9%",
          "--foreground": "0 0% 98%",
          "--card": "0 0% 3.9%",
          "--card-foreground": "0 0% 98%",
          "--popover": "0 0% 3.9%",
          "--popover-foreground": "0 0% 98%",
         " --primary": "0 0% 98%",
         "--primary-foreground": "240 5.9% 10%",
          "--secondary": "207 90% 54%", // Changed to a blue color (approximate)
          "--secondary-foreground": "0 0% 98%",
          "--muted": "0 0% 14.9%",
          "--muted-foreground": "0 0% 63.9%",
          "--accent": "0 0% 8%",
          "--accent-foreground": "0 0% 98%",
          "--destructive": "0 62.8% 30.6%",
          "--destructive-foreground": "0 0% 98%",
          "--border": "0 0% 14.9%",
          "--input": "0 0% 14.9%",
          "--ring": "0 0% 83.1%",
			},
      }
    })
  ],
  content: [
    './**/*.{templ,html}', 
    "./**/*.html", 
    "./**/*.templ", 
    "./**/*.go",
  ],
  safelist: [
    {
      pattern: /^uk-/
    }
  ],
  theme: {
    extend: {
      colors: {
        border: "hsl(var(--border) / <alpha-value>)",
				input: "hsl(var(--input) / <alpha-value>)",
				ring: "hsl(var(--ring) / <alpha-value>)",
				background: "hsl(var(--background) / <alpha-value>)",
				foreground: "hsl(var(--foreground) / <alpha-value>)",
				primary: {
					DEFAULT: "hsl(var(--primary) / <alpha-value>)",
					foreground: "hsl(var(--primary-foreground) / <alpha-value>)"
				},
				secondary: {
					DEFAULT: "hsl(var(--secondary) / <alpha-value>)",
					foreground: "hsl(var(--secondary-foreground) / <alpha-value>)"
				},
				destructive: {
					DEFAULT: "hsl(var(--destructive) / <alpha-value>)",
					foreground: "hsl(var(--destructive-foreground) / <alpha-value>)"
				},
				muted: {
					DEFAULT: "hsl(var(--muted) / <alpha-value>)",
					foreground: "hsl(var(--muted-foreground) / <alpha-value>)"
				},
				accent: {
					DEFAULT: "hsl(var(--accent) / <alpha-value>)",
					foreground: "hsl(var(--accent-foreground) / <alpha-value>)"
				},
				popover: {
					DEFAULT: "hsl(var(--popover) / <alpha-value>)",
					foreground: "hsl(var(--popover-foreground) / <alpha-value>)"
				},
				card: {
					DEFAULT: "hsl(var(--card) / <alpha-value>)",
					foreground: "hsl(var(--card-foreground) / <alpha-value>)"
				}
			},
			borderRadius: {
				lg: "var(--radius)",
				md: "calc(var(--radius) - 2px)",
				sm: "calc(var(--radius) - 4px)"
			},
      keyframes: {
        "accordion-down": {
          from: { height: 0 },
          to: { height: "var(--radix-accordion-content-height)" },
        },
        "accordion-up": {
          from: { height: "var(--radix-accordion-content-height)" },
          to: { height: 0 },
        },
      },
      animation: {
        "accordion-down": "accordion-down 0.2s ease-out",
        "accordion-up": "accordion-up 0.2s ease-out",
      },
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
};