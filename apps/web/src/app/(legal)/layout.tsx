import Link from "next/link";

export default function LegalLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const legalNav = [
    { name: "Términos y Condiciones", href: "/terminos" },
    { name: "Política de Privacidad", href: "/privacidad" },
    { name: "Política de Cookies", href: "/cookies" },
    { name: "Uso Aceptable", href: "/uso-aceptable" },
    { name: "Reembolsos y Cancelación", href: "/reembolsos" },
    { name: "Tratamiento de Datos (DPA)", href: "/dpa" },
  ];

  return (
    <div className="pt-40 pb-24 px-4 bg-background min-h-screen">
      <div className="max-w-6xl mx-auto grid grid-cols-1 md:grid-cols-12 gap-12 md:gap-24">
        
        {/* Sidebar Sticky */}
        <aside className="md:col-span-3">
          <div className="sticky top-32">
            <h2 className="text-heading font-display text-foreground mb-6">Legal</h2>
            <nav className="flex flex-col gap-3">
              {legalNav.map((item) => (
                <Link 
                  key={item.href} 
                  href={item.href}
                  className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
                >
                  {item.name}
                </Link>
              ))}
            </nav>
          </div>
        </aside>

        {/* Contenido Legal */}
        <div className="md:col-span-9 bg-white">
          <div className="prose prose-neutral max-w-none prose-headings:font-display prose-headings:font-medium prose-headings:tracking-tight prose-p:text-muted-foreground prose-li:text-muted-foreground prose-a:text-foreground prose-a:underline hover:prose-a:text-muted-foreground">
            {children}
          </div>
        </div>

      </div>
    </div>
  );
}
