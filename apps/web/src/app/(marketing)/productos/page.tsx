import Link from "next/link";
import { Button } from "@/components/ui/button";
import { LayoutTemplate, CreditCard, BarChart3, Package, Zap, Shield, Globe, Smartphone } from "lucide-react";

export default function ProductosPage() {
  const features = [
    { icon: <LayoutTemplate className="w-8 h-8 text-foreground" />, title: "Editor Visual Puck", desc: "Arrastra, suelta y publica. Tu tienda lista en minutos sin tocar una línea de código. Plantillas diseñadas para convertir." },
    { icon: <CreditCard className="w-8 h-8 text-foreground" />, title: "Pagos Nativos", desc: "Conectado con Pago Móvil, C2P, Zelle, Binance y las principales pasarelas de Venezuela. Cobra en bolívares y dólares." },
    { icon: <Package className="w-8 h-8 text-foreground" />, title: "Dropshipping Integrado", desc: "Accede a nuestra red de proveedores verificados. Publica productos y nosotros nos encargamos de gestionar el despacho." },
    { icon: <BarChart3 className="w-8 h-8 text-foreground" />, title: "Analíticas en Tiempo Real", desc: "Mide visitas, conversión y ventas. Toma decisiones basadas en datos con nuestro dashboard conectado a tu tienda." },
  ];

  const steps = [
    { title: "Crea tu cuenta", desc: "Regístrate en minutos. No necesitas tarjeta de crédito para empezar." },
    { title: "Diseña tu tienda", desc: "Elige una plantilla, personaliza colores y sube tu logo con el editor visual." },
    { title: "Conecta pagos", desc: "Activa Pago Móvil, Zelle o Binance Pay con un solo clic." },
    { title: "Publica y vende", desc: "Tu tienda estará lista para recibir pedidos en cuestión de minutos." },
  ];

  const integrations = [
    { icon: <Shield className="w-6 h-6" />, name: "Wayu Pay" },
    { icon: <CreditCard className="w-6 h-6" />, name: "VeriPay" },
    { icon: <Globe className="w-6 h-6" />, name: "VEXPay" },
    { icon: <Smartphone className="w-6 h-6" />, name: "Payku" },
    { icon: <Zap className="w-6 h-6" />, name: "Pabilo" },
  ];

  return (
    <div className="pt-40 pb-24 px-4 bg-background min-h-screen">
      <div className="max-w-7xl mx-auto">
        <div className="max-w-3xl mb-20">
          <h1 className="text-display-1 mb-6 text-foreground tracking-tight">Todo lo que necesitas para vender.</h1>
          <p className="text-body-lg text-muted-foreground">Una plataforma completa. Sin integraciones complejas, sin plugins de terceros.</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-x-16 gap-y-16 mb-32">
          {features.map((feature, idx) => (
            <div key={idx} className="flex flex-col gap-6">
              <div className="w-16 h-16 bg-neutral-50 border border-border rounded-2xl flex items-center justify-center">{feature.icon}</div>
              <div>
                <h3 className="text-heading-lg font-display text-foreground mb-3">{feature.title}</h3>
                <p className="text-body text-muted-foreground max-w-md">{feature.desc}</p>
              </div>
            </div>
          ))}
        </div>

        <div className="border-t border-border pt-20 mb-32">
          <h2 className="text-display-2 mb-12">Cómo funciona en 4 pasos.</h2>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
            {steps.map((step, idx) => (
              <div key={idx}>
                <div className="text-display-3 font-display text-accent mb-4">0{idx + 1}</div>
                <h3 className="text-heading font-display text-foreground mb-2">{step.title}</h3>
                <p className="text-body-sm text-muted-foreground">{step.desc}</p>
              </div>
            ))}
          </div>
        </div>

        <div className="border-t border-border pt-20 mb-32">
          <h2 className="text-display-2 mb-8">Integraciones nativas.</h2>
          <p className="text-body-lg text-muted-foreground mb-12 max-w-2xl">Trabajamos con las pasarelas de pago que ya usas en Venezuela. Sin configuraciones raras, sin complicaciones.</p>
          <div className="flex flex-wrap gap-6">
            {integrations.map((int, idx) => (
              <div key={idx} className="flex items-center gap-3 px-6 py-4 bg-neutral-50 rounded-full border border-border">
                {int.icon}
                <span className="text-body font-medium">{int.name}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="border-t border-border pt-16 flex flex-col md:flex-row items-start md:items-center justify-between gap-8">
          <h2 className="text-display-2 text-foreground max-w-xl">Empieza a construir hoy.</h2>
          <Link href="/registro"><Button variant="default" size="lg">Crear mi tienda gratis</Button></Link>
        </div>
      </div>
    </div>
  );
}
