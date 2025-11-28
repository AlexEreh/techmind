import "@/styles/globals.css";
import { Metadata, Viewport } from "next";
import clsx from "clsx";

import { Providers } from "./providers";

import { siteConfig } from "@/config/site";
import { fontSans } from "@/config/fonts";

export const metadata: Metadata = {
  title: {
    default: "DataRush - Менеджер файлов",
    template: `%s - DataRush`,
  },
  description: "Менеджер файлов для организаций",
  icons: {
    icon: "/favicon.ico",
  },
};

export const viewport: Viewport = {
  themeColor: [
    { media: "(prefers-color-scheme: light)", color: "white" },
    { media: "(prefers-color-scheme: dark)", color: "black" },
  ],
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html suppressHydrationWarning lang="ru">
      <head />
      <body
          className={clsx(
              "min-h-screen text-foreground font-sans antialiased",
              fontSans.variable,
          )}
          style={{
              backgroundImage:
                  'url("https://media.sketchfab.com/models/c46b152c6f904de799bff4415794d07c/thumbnails/63c01b440f1848f588b99edd9a47e408/53d062c22a204cdb848f10363b21efa1.jpeg")',
              backgroundSize: "cover",
              backgroundPosition: "center",
              backgroundRepeat: "no-repeat",
          }}
      >
        <Providers themeProps={{ attribute: "class", defaultTheme: "dark" }}>
          {children}
        </Providers>
      </body>
    </html>
  );
}
