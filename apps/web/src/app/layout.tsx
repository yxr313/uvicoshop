import type { Metadata } from "next";
import { Bricolage_Grotesque, Instrument_Sans } from "next/font/google";
import { Toaster } from "sonner";
import "./globals.css";

const bricolage = Bricolage_Grotesque({
  subsets: ["latin"],
  variable: "--font-bricolage",
  display: "swap",
});

const instrument = Instrument_Sans({
  subsets: ["latin"],
  variable: "--font-instrument",
  display: "swap",
});

export const metadata: Metadata = {
  title: "Uvico — Crea tu tienda online en Venezuela",
  description: "El Ecommerce Builder que impulsa a los emprendedores venezolanos.",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="es-VE" className={`${bricolage.variable} ${instrument.variable}`}>
      <body className="font-body antialiased bg-background text-foreground min-h-dvh">
        {children}
        <Toaster position="top-center" richColors />
      </body>
    </html>
  );
}
