import Link from "next/link";
import { ChromaticImage } from "@/components/ui/chromatic-image";

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex h-[100dvh] w-full bg-background overflow-hidden">
      
      {/* Lado Izquierdo: Imagen a pantalla completa (Sin texto) */}
      <div className="hidden lg:block w-1/2 relative bg-[#0a0a0a] overflow-hidden border-r border-border/10">
        <ChromaticImage
          src="https://assets.aceternity.com/screenshots/8131a006-884a-4444-85a6-aee9d56af136.webp"
          alt="Uvico Commerce"
          className="w-full h-full"
        />
      </div>

      {/* Lado Derecho: Formulario */}
      <div className="flex-1 flex flex-col justify-center px-6 sm:px-12 md:px-24 h-full overflow-y-auto bg-background">
        <div className="w-full max-w-[400px] mx-auto py-12">
          {/* Logo UVICO arriba del formulario */}
          <Link href="/" className="font-display font-bold text-2xl tracking-tight text-foreground inline-block mb-12">
            UVICO
          </Link>
          {children}
        </div>
      </div>

    </div>
  );
}
