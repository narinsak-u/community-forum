import type { Metadata } from "next";
import { JetBrains_Mono, Space_Grotesk } from "next/font/google";
import { Providers } from "@/components/providers";
import "./globals.css";

const jetbrainsMono = JetBrains_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
  display: "swap",
});

const spaceGrotesk = Space_Grotesk({
  subsets: ["latin"],
  variable: "--font-display",
  display: "swap",
});

export const metadata: Metadata = {
  title: "Midnight Forge - Terminal Forum for Architects & Engineers",
  description:
    "Midnight Forge is a terminal-style technical forum for architects, engineers, and protocol designers. Discuss systems, share docs, build the network.",
  authors: [{ name: "Midnight Forge" }],
  openGraph: {
    title: "Midnight Forge - Terminal Forum",
    description: "Terminal-style technical forum for architects and engineers.",
    type: "website",
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html
      lang="en"
      className={`${jetbrainsMono.variable} ${spaceGrotesk.variable} mx-auto w-full`}
      suppressHydrationWarning
    >
      <body suppressHydrationWarning>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
