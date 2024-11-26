import type { Metadata } from "next";
import localFont from "next/font/local";
import "./globals.css";
import Providers from "./providers";

const oldschool = localFont({
  src: "./fonts/OldschoolGroteskCompact-ExtraBold.woff2",
  variable: "--oldschool",
  weight: "800",
});

const monzoSansDisplay = localFont({
  src: [
    {
      path: "./fonts/MonzoSansDisplay-Regular.woff2",
      weight: "400",
    },
  ],
  variable: "--monzosansdisplay",
});

const monzoSansText = localFont({
  src: [
    {
      path: "./fonts/MonzoSansText-Italic.woff2",
      weight: "400",
      style: "italic"
    },
    {
      path: "./fonts/MonzoSansText-Regular.woff2",
      weight: "500",
    },
    {
      path: "./fonts/MonzoSansText-SemiBold.woff2",
      weight: "600",
    },
    {
      path: "./fonts/MonzoSansText-Bold.woff2",
      weight: "700",
    },
  ],
  variable: "--monzosanstext",
});


export const metadata: Metadata = {
  title: "monzopanel",
  description: "creating an analytics pipeline with Monzo's core technologies",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${oldschool.variable} ${monzoSansDisplay.variable} ${monzoSansText.variable}`}>
        <Providers>
          {children}
        </Providers>
      </body>
    </html>
  );
}
