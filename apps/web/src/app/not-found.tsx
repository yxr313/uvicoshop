import Link from "next/link";
import { ArrowLeft, Search } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function NotFound() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-8 bg-background">
      <div className="max-w-md w-full text-center">
        <p className="text-[6rem] leading-none font-display font-medium tracking-tighter text-neutral-200 mb-4">
          404
        </p>
        <h1 className="text-display-3 font-display text-foreground mb-4">
          Página no encontrada.
        </h1>
        <p className="text-body text-muted-foreground mb-10">
          La página que buscas no existe, ha sido movida o está en mantenimiento.
        </p>
        
        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          <Link href="/" className="w-full sm:w-auto">
            <Button className="w-full rounded-full bg-accent text-accent-foreground hover:bg-accent-hover">
              <ArrowLeft className="mr-2" size={18} /> Volver al inicio
            </Button>
          </Link>
          <Link href="/ayuda" className="w-full sm:w-auto">
            <Button variant="outline" className="w-full rounded-full">
              <Search className="mr-2" size={18} /> Centro de ayuda
            </Button>
          </Link>
        </div>
      </div>
    </div>
  );
}
