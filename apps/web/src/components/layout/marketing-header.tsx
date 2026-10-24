import Link from "next/link";
import { Button } from "@/components/ui/button";

export function MarketingHeader() {
  const navItems = [
    { name: "Productos", href: "/productos" },
    { name: "Proveedores", href: "/proveedores" },
    { name: "Precios", href: "/planes" },
    { name: "Recursos", href: "/recursos" },
  ];

  return (
    <header className="fixed top-6 left-0 right-0 z-50 px-4 flex justify-center w-full">
      <div className="flex items-center justify-between w-full max-w-4xl px-6 py-3 bg-white/80 backdrop-blur-md border border-border rounded-full shadow-sm">
        {/* Logo */}
        <Link href="/" className="font-display font-bold text-xl tracking-tight text-foreground">
          UVICO
        </Link>

        {/* Enlaces centrados */}
        <nav className="hidden md:flex items-center gap-8">
          {navItems.map((item) => (
            <Link
              key={item.name}
              href={item.href}
              className="relative text-sm font-medium text-muted-foreground hover:text-foreground transition-colors py-1
                         after:absolute after:bottom-0 after:left-1/2 after:-translate-x-1/2 after:h-[1px] after:w-0 after:bg-foreground after:transition-all after:duration-300 hover:after:w-full"
            >
              {item.name}
            </Link>
          ))}
        </nav>

        {/* Botones a la derecha */}
        <div className="flex items-center gap-2">
          <Link href="/login">
            <Button variant="ghost" size="sm" className="hidden sm:inline-flex">
              Iniciar sesión
            </Button>
          </Link>
          <Link href="/registro">
            <Button variant="default" size="sm">
              Crear tienda
            </Button>
          </Link>
        </div>
      </div>
    </header>
  );
}
