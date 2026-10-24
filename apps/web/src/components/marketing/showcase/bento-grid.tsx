import { LayoutTemplate, Image as ImageIcon, Type, Settings, GripVertical, ArrowUpRight, ShoppingBag } from "lucide-react";

export function BentoGrid() {
  return (
    <section className="py-24 px-4 bg-background">
      <div className="max-w-6xl mx-auto">
        
        {/* H2 de la sección */}
        <h2 className="text-display-2 mb-16 text-foreground text-center max-w-3xl mx-auto">
          Un solo lugar para crear, gestionar y vender.
        </h2>

        {/* Grid de 2 columnas (Bento) */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          
          {/* Bloque 1: El Editor Visual */}
          <div className="bg-neutral-50 border border-border rounded-[2rem] p-8 md:p-10 flex flex-col justify-between overflow-hidden relative min-h-[450px] group">
            <div className="relative z-10 mb-8">
              <h3 className="text-heading-lg mb-3">Tu tienda, tu diseño.</h3>
              <p className="text-body text-muted-foreground max-w-sm">
                Arrastra, suelta y publica. Un editor visual diseñado para que tu marca destaque sin tocar una línea de código.
              </p>
            </div>
            
            {/* Mockup Visual: Editor Drag & Drop */}
            <div className="relative w-full h-[240px] bg-white border border-border rounded-t-2xl shadow-2xl translate-y-6 group-hover:translate-y-2 transition-transform duration-500 flex">
              {/* Sidebar del editor */}
              <div className="w-16 border-r border-border flex flex-col items-center py-4 gap-6 bg-neutral-50/50">
                <div className="w-8 h-8 rounded bg-neutral-100 flex items-center justify-center text-neutral-500"><LayoutTemplate size={16} /></div>
                <div className="w-8 h-8 rounded flex items-center justify-center text-neutral-400 hover:text-foreground transition-colors"><ImageIcon size={16} /></div>
                <div className="w-8 h-8 rounded flex items-center justify-center text-neutral-400 hover:text-foreground transition-colors"><Type size={16} /></div>
                <div className="w-8 h-8 rounded flex items-center justify-center text-neutral-400 hover:text-foreground transition-colors"><Settings size={16} /></div>
              </div>
              
              {/* Lienzo (Canvas) */}
              <div className="flex-1 p-6 relative bg-[radial-gradient(#e5e5e5_1px,transparent_1px)] [background-size:16px_16px]">
                {/* Bloque seleccionado (simulando Puck) */}
                <div className="relative border-2 border-accent border-dashed rounded-lg p-4 bg-white/80 backdrop-blur-sm">
                  {/* Etiqueta del bloque */}
                  <div className="absolute -top-3 -left-2 bg-accent text-accent-foreground text-xs font-semibold px-2 py-0.5 rounded flex items-center gap-1 shadow-sm">
                    <GripVertical size={12} />
                    Hero Section
                  </div>
                  
                  <div className="w-3/4 h-6 bg-foreground rounded mb-2" />
                  <div className="w-1/2 h-3 bg-neutral-200 rounded mb-4" />
                  <div className="w-24 h-8 bg-accent rounded-full" />
                </div>
              </div>
            </div>
          </div>

          {/* Bloque 2: Pagos y Dropshipping (Notificaciones Apiladas) */}
          <div className="bg-neutral-50 border border-border rounded-[2rem] p-8 md:p-10 flex flex-col justify-between overflow-hidden relative min-h-[450px]">
            <div className="relative z-10 mb-4">
              <h3 className="text-heading-lg mb-3">Cobra en bolívares y dólares.</h3>
              <p className="text-body text-muted-foreground max-w-sm">
                Conectado de forma nativa con las pasarelas venezolanas y con nuestra red de proveedores verificados.
              </p>
            </div>
            
            {/* Mockup Visual: Dashboard de Ventas y Notificaciones */}
            <div className="relative w-full h-[260px] mt-auto">
              
              {/* Gráfico de fondo (Curva SVG) */}
              <div className="absolute inset-x-0 bottom-0 w-full h-[160px]">
                <svg viewBox="0 0 400 150" className="w-full h-full" preserveAspectRatio="none">
                  <defs>
                    <linearGradient id="gradient" x1="0" x2="0" y1="0" y2="1">
                      <stop offset="0%" stopColor="#d4d700" stopOpacity="0.3" />
                      <stop offset="100%" stopColor="#d4d700" stopOpacity="0" />
                    </linearGradient>
                  </defs>
                  <path d="M0,150 C50,150 100,100 150,80 C200,60 250,20 300,30 C350,40 400,10 400,10 L400,150 Z" fill="url(#gradient)" />
                  <path d="M0,150 C50,150 100,100 150,80 C200,60 250,20 300,30 C350,40 400,10 400,10" fill="none" stroke="#d4d700" strokeWidth="3" strokeLinecap="round" />
                </svg>
              </div>

              {/* Notificaciones Apiladas (Stacked Cards) */}
              <div className="absolute inset-x-0 bottom-6 flex flex-col gap-3 px-2">
                
                {/* Venta 3 (Fondo) */}
                <div className="w-[88%] mx-auto bg-white/60 border border-border p-3.5 rounded-xl shadow-sm flex items-center justify-between translate-y-[-15px] scale-90 opacity-40 backdrop-blur-md">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 bg-neutral-100 rounded-full flex items-center justify-center text-neutral-500">
                      <ShoppingBag size={14} />
                    </div>
                    <div>
                      <p className="text-xs font-semibold text-foreground">Tech Store</p>
                      <p className="text-[10px] text-muted-foreground">Pago Móvil</p>
                    </div>
                  </div>
                  <p className="text-sm font-display font-bold text-foreground">+$15.00</p>
                </div>

                {/* Venta 2 (Medio) */}
                <div className="w-[94%] mx-auto bg-white border border-border p-3.5 rounded-xl shadow-lg flex items-center justify-between translate-y-[-5px] scale-95 opacity-80 z-10">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 bg-neutral-100 rounded-full flex items-center justify-center text-neutral-500">
                      <ShoppingBag size={14} />
                    </div>
                    <div>
                      <p className="text-xs font-semibold text-foreground">Zapatos Caracas</p>
                      <p className="text-[10px] text-muted-foreground">Wayu Pay</p>
                    </div>
                  </div>
                  <p className="text-sm font-display font-bold text-foreground">+$45.50</p>
                </div>

                {/* Venta 1 (Frente - Destacada) */}
                <div className="w-full bg-white border border-border p-4 rounded-xl shadow-xl flex items-center justify-between relative z-20">
                  <div className="absolute -top-2.5 -right-2.5 w-6 h-6 bg-accent text-accent-foreground rounded-full flex items-center justify-center shadow-sm">
                    <ArrowUpRight size={14} strokeWidth={3} />
                  </div>
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-accent-muted rounded-full flex items-center justify-center text-foreground">
                      <ShoppingBag size={18} />
                    </div>
                    <div>
                      <p className="text-sm font-semibold text-foreground">Bodega El Sol</p>
                      <p className="text-xs text-muted-foreground">Pago Móvil validado</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="text-lg font-display font-bold text-foreground">+$120.00</p>
                    <p className="text-xs font-medium text-foreground/60">Hace un momento</p>
                  </div>
                </div>

              </div>
            </div>
          </div>

        </div>
      </div>
    </section>
  );
}
