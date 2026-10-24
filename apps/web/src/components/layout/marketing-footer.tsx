import Link from "next/link";
import { Instagram, Github, ArrowUpRight } from "lucide-react";

export function MarketingFooter() {
  const sections = [
    {
      title: "Producto",
      links: [
        { name: "Editor visual", href: "/productos" },
        { name: "Dropshipping", href: "/proveedores" },
        { name: "Precios", href: "/planes" },
        { name: "Proveedores", href: "/proveedores" },
      ],
    },
    {
      title: "Compañía",
      links: [
        { name: "Nosotros", href: "/nosotros" },
        { name: "Blog", href: "/blog" },
        { name: "Yamal Network", href: "/yamal-network" },
        { name: "Trabajos", href: "/trabajos" },
        { name: "Contacto", href: "/contacto" },
      ],
    },
    {
      title: "Soporte",
      links: [
        { name: "Centro de ayuda", href: "/ayuda" },
        { name: "Guías", href: "/recursos" },
        { name: "Estado del sistema", href: "/status" },
      ],
    },
    {
      title: "Contacto",
      links: [
        { name: "contacto@uvico.shop", href: "mailto:contacto@uvico.shop" },
        { name: "privacidad@uvico.shop", href: "mailto:privacidad@uvico.shop" },
        { name: "yxr@uvico.shop", href: "mailto:yxr@uvico.shop" },
        { name: "yxr.313@pm.me", href: "mailto:yxr.313@pm.me" },
      ],
    },
  ];

  const legalLinks = [
    { name: "Términos y Condiciones", href: "/terminos" },
    { name: "Política de Privacidad", href: "/privacidad" },
    { name: "Cookies", href: "/cookies" },
    { name: "Uso Aceptable", href: "/uso-aceptable" },
    { name: "Reembolsos", href: "/reembolsos" },
    { name: "DPA", href: "/dpa" },
  ];

  const XIcon = () => (
    <svg viewBox="0 0 24 24" fill="currentColor" className="w-5 h-5">
      <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z" />
    </svg>
  );

  const WhatsAppIcon = () => (
    <svg viewBox="0 0 24 24" fill="currentColor" className="w-5 h-5">
      <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L0 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z" />
    </svg>
  );

  return (
    <footer className="py-24 px-4 bg-background border-t border-border">
      <div className="max-w-7xl mx-auto">
        
        <div className="grid grid-cols-2 gap-x-8 gap-y-12 lg:grid-cols-6 mb-20">
          <div className="col-span-2">
            <Link href="/" className="font-display font-bold text-2xl tracking-tight text-foreground inline-block mb-4">
              UVICO
            </Link>
            <p className="text-body text-muted-foreground max-w-xs mb-8">
              La plataforma todo en uno para crear tu tienda, gestionar dropshipping y cobrar en Venezuela.
            </p>
            
            {/* Redes Sociales y Canales de WhatsApp Alineados */}
            <div className="flex items-center gap-5">
              <a href="https://instagram.com/uvico.ve" target="_blank" rel="noreferrer" className="text-muted-foreground hover:text-foreground transition-colors" aria-label="Instagram Uvico"><Instagram className="w-5 h-5" /></a>
              <a href="https://x.com/uvicoshop" target="_blank" rel="noreferrer" className="text-muted-foreground hover:text-foreground transition-colors" aria-label="X Uvico"><XIcon /></a>
              <a href="https://github.com/yxr313" target="_blank" rel="noreferrer" className="text-muted-foreground hover:text-foreground transition-colors" aria-label="GitHub Uvico"><Github className="w-5 h-5" /></a>
              
              {/* WhatsApp de Uvico */}
              <a 
                href="https://whatsapp.com/channel/0029VbEDp1CICVfc8CaOi02w" 
                target="_blank" 
                rel="noreferrer" 
                className="text-muted-foreground hover:text-foreground transition-colors" 
                aria-label="Canal de WhatsApp de Uvico"
                title="Canal de WhatsApp de Uvico"
              >
                <WhatsAppIcon />
              </a>

              {/* WhatsApp de Yamal */}
              <a 
                href="https://whatsapp.com/channel/0029Vb8t14cKgsO0pMXKGg0L" 
                target="_blank" 
                rel="noreferrer" 
                className="text-muted-foreground hover:text-foreground transition-colors" 
                aria-label="Canal de WhatsApp de Yamal Network"
                title="Canal de WhatsApp de Yamal Network"
              >
                <WhatsAppIcon />
              </a>
            </div>
          </div>

          {sections.map((section, idx) => (
            <div key={idx}>
              <h3 className="text-heading font-display text-foreground mb-6">{section.title}</h3>
              <ul className="space-y-4">
                {section.links.map((link, linkIdx) => (
                  <li key={linkIdx}>
                    <Link href={link.href} className="group inline-flex items-center text-body-sm font-medium text-muted-foreground hover:text-foreground transition-colors">
                      {link.name}
                      <ArrowUpRight className="w-3.5 h-3.5 ml-1 opacity-0 -translate-x-2 group-hover:opacity-100 group-hover:translate-x-0 transition-all duration-300" />
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        <div className="pt-8 border-t border-border flex flex-col md:flex-row justify-between items-center gap-6">
          <p className="text-sm font-medium text-muted-foreground">
            &copy; {new Date().getFullYear()} Uvico, S.A. Todos los derechos reservados.
          </p>
          <ul className="flex flex-wrap justify-center gap-x-8 gap-y-3">
            {legalLinks.map((link, idx) => (
              <li key={idx}>
                <Link href={link.href} className="group inline-flex items-center text-sm font-medium text-muted-foreground underline underline-offset-4 hover:text-foreground transition-colors">
                  {link.name}
                  <ArrowUpRight className="w-3.5 h-3.5 ml-1 opacity-0 -translate-x-2 group-hover:opacity-100 group-hover:translate-x-0 transition-all duration-300" />
                </Link>
              </li>
            ))}
          </ul>
        </div>

      </div>
    </footer>
  );
}
