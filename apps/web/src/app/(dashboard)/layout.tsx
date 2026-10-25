"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { sidebarData } from "@/config/sidebar";
import { Button } from "@/components/ui/button";
import { PanelLeftClose, PanelLeftOpen, Bell, Package, ShoppingCart } from "lucide-react";
import { cn } from "@/lib/utils";
import "@puckeditor/core/puck.css"; // Importante para el editor

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [showNotifications, setShowNotifications] = useState(false);
  const pathname = usePathname();

  const getBreadcrumbs = () => {
    if (pathname === "/dashboard") return "Overview";
    const segments = pathname.split("/").filter(Boolean);
    const section = segments[1];
    const subSection = segments[2];

    if (section === "ecommerce" && subSection) {
      const map: Record<string, string> = { products: "Productos", orders: "Pedidos", customers: "Clientes" };
      return `E-commerce / ${map[subSection] || subSection}`;
    }
    if (section === "workspaces" && subSection) {
      const map: Record<string, string> = { projects: "Proyectos", tasks: "Tareas", calendar: "Calendario" };
      return `Workspaces / ${map[subSection] || subSection}`;
    }
    if (section === "account" && subSection) {
      const map: Record<string, string> = { profile: "Perfil", security: "Seguridad", notifications: "Notificaciones" };
      return `Account / ${map[subSection] || subSection}`;
    }
    if (section === "editor") return "Editor Visual";
    return section.charAt(0).toUpperCase() + section.slice(1);
  };

  return (
    <div className="flex min-h-screen bg-neutral-50">
      {/* Sidebar (código idéntico al anterior) */}
      <aside className={cn("bg-white border-r border-border hidden md:flex flex-col fixed inset-y-0 z-50 transition-all duration-300", isCollapsed ? "w-20" : "w-64")}>
        <div className={cn("h-16 flex items-center border-b border-border px-6", isCollapsed ? "justify-center px-0" : "justify-between")}>
          {!isCollapsed && <Link href="/" className="font-display font-bold text-xl tracking-tight text-foreground">UVICO<span className="text-accent">.</span></Link>}
          <button onClick={() => setIsCollapsed(!isCollapsed)} className="p-1.5 rounded-md hover:bg-neutral-100 text-muted-foreground transition-colors">
            {isCollapsed ? <PanelLeftOpen size={20} /> : <PanelLeftClose size={20} />}
          </button>
        </div>
        <div className="flex-1 overflow-y-auto py-6 px-3 space-y-8 no-scrollbar">
          {sidebarData.map((section, idx) => (
            <div key={idx}>
              {!isCollapsed && <h4 className="text-[11px] font-semibold text-neutral-400 uppercase tracking-wider mb-3 px-3">{section.title}</h4>}
              <nav className="space-y-1">
                {section.items.map((item, i) => (
                  <div key={i}>
                    {item.items ? (
                      <div className="space-y-1">
                        <div className={cn("flex items-center gap-3 px-3 py-2 text-sm font-medium text-foreground", isCollapsed && "justify-center px-0")}>
                          <item.icon size={18} className="text-neutral-500 shrink-0" />
                          {!isCollapsed && item.title}
                        </div>
                        {!isCollapsed && (
                          <div className="pl-9 space-y-1">
                            {item.items.map((sub, j) => (
                              <Link key={j} href={sub.url} className={cn("flex items-center gap-2 px-3 py-1.5 text-sm rounded-md transition-colors", pathname === sub.url ? "bg-neutral-100 text-foreground font-medium" : "text-muted-foreground hover:text-foreground hover:bg-neutral-100")}>
                                {sub.icon && <sub.icon size={14} className="opacity-60 shrink-0" />}
                                {sub.title}
                              </Link>
                            ))}
                          </div>
                        )}
                      </div>
                    ) : (
                      <Link href={item.url!} className={cn("flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-md transition-colors", pathname === item.url ? "bg-neutral-100 text-foreground" : "text-muted-foreground hover:text-foreground hover:bg-neutral-100", isCollapsed && "justify-center px-0")} title={isCollapsed ? item.title : undefined}>
                        <item.icon size={18} className="shrink-0" />
                        {!isCollapsed && item.title}
                      </Link>
                    )}
                  </div>
                ))}
              </nav>
            </div>
          ))}
        </div>
      </aside>

      <main className={cn("flex-1 flex flex-col min-h-screen transition-all duration-300", isCollapsed ? "md:ml-20" : "md:ml-64")}>
        <header className="h-16 bg-white border-b border-border flex items-center justify-between px-8 sticky top-0 z-40">
          <div className="flex items-center gap-2 text-sm">
            <span className="text-muted-foreground">Dashboard</span>
            <span className="text-neutral-300">/</span>
            <span className="font-medium text-foreground">{getBreadcrumbs()}</span>
          </div>
          <div className="flex items-center gap-4">
            <div className="relative">
              <button onClick={() => setShowNotifications(!showNotifications)} className="relative p-2 rounded-full hover:bg-neutral-100 text-muted-foreground transition-colors">
                <Bell size={20} />
                <span className="absolute top-1.5 right-1.5 w-2 h-2 bg-red-500 rounded-full border-2 border-white" />
              </button>
              {showNotifications && (
                <div className="absolute right-0 top-12 w-80 bg-white border border-border rounded-xl shadow-xl z-50 overflow-hidden animate-in fade-in slide-in-from-top-2">
                  <div className="p-4 border-b border-border flex justify-between items-center bg-neutral-50">
                    <span className="font-medium text-sm text-foreground">Notificaciones</span>
                    <span className="text-xs text-accent font-semibold cursor-pointer">Marcar leídas</span>
                  </div>
                  <div className="divide-y divide-border">
                    <div className="p-4 flex gap-3 hover:bg-neutral-50 transition-colors cursor-pointer">
                      <div className="w-8 h-8 rounded-full bg-green-100 text-green-600 flex items-center justify-center shrink-0"><ShoppingCart size={14} /></div>
                      <div><p className="text-sm font-medium text-foreground">Nueva venta realizada</p><p className="text-xs text-muted-foreground mt-0.5">Pedido #104 por $45.00</p><p className="text-[10px] text-neutral-400 mt-1">Hace 2 minutos</p></div>
                    </div>
                    <div className="p-4 flex gap-3 hover:bg-neutral-50 transition-colors cursor-pointer">
                      <div className="w-8 h-8 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center shrink-0"><Package size={14} /></div>
                      <div><p className="text-sm font-medium text-foreground">Stock bajo</p><p className="text-xs text-muted-foreground mt-0.5">El producto "Camisa Oversize" tiene 2 unidades.</p><p className="text-[10px] text-neutral-400 mt-1">Hace 1 hora</p></div>
                    </div>
                  </div>
                  <div className="p-2 border-t border-border text-center bg-neutral-50">
                    <button className="text-xs font-medium text-muted-foreground hover:text-foreground w-full py-1">Ver todas</button>
                  </div>
                </div>
              )}
            </div>
            <Button variant="outline" size="sm" asChild><Link href="/">Ver mi tienda</Link></Button>
            <div className="w-8 h-8 rounded-full bg-accent text-accent-foreground flex items-center justify-center font-bold text-xs">YX</div>
          </div>
        </header>
        <div className="flex-1 p-8">{children}</div>
      </main>
    </div>
  );
}
