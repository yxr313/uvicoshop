import Link from "next/link";
import { Search, BookOpen, CreditCard, Boxes, ShieldCheck } from "lucide-react";

export default function AyudaPage() {
  const categories = [
    { icon: <BookOpen />, title: "Primeros pasos", desc: "Aprende a configurar tu tienda desde cero." },
    { icon: <CreditCard />, title: "Pagos y Pasarelas", desc: "Configuración de Pago Móvil, Zelle y más." },
    { icon: <Boxes />, title: "Dropshipping", desc: "Cómo conectar proveedores y gestionar pedidos." },
    { icon: <ShieldCheck />, title: "Seguridad", desc: "Protección de cuenta y retención de fondos." },
  ];

  return (
    <div className="pt-40 pb-24 px-4 bg-background min-h-screen">
      <div className="max-w-4xl mx-auto">
        <div className="text-center mb-16">
          <h1 className="text-display-1 mb-6 text-foreground tracking-tight">Centro de ayuda.</h1>
          <div className="relative max-w-xl mx-auto">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-muted-foreground w-5 h-5" />
            <input 
              type="text" 
              placeholder="Busca un artículo..." 
              className="w-full h-14 pl-12 pr-4 rounded-full border border-border bg-white focus:outline-none focus:border-neutral-400 transition-all text-body shadow-sm"
            />
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {categories.map((cat, idx) => (
            <Link key={idx} href="/recursos" className="p-8 border border-border rounded-2xl hover:bg-neutral-50 transition-colors group">
              <div className="w-12 h-12 bg-white border border-border rounded-xl flex items-center justify-center mb-6 text-foreground group-hover:border-neutral-300">
                {cat.icon}
              </div>
              <h3 className="text-heading font-display text-foreground mb-2">{cat.title}</h3>
              <p className="text-body-sm text-muted-foreground">{cat.desc}</p>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}
