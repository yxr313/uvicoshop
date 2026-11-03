"use client";

import { useState } from "react";
import Link from "next/link";
import { IconArrowRight, IconTrash, IconCheck, IconBell } from "@tabler/icons-react";

import { PageHeader } from "@/components/page-header";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { useNotifications, timeAgo } from "@/hooks/use-notifications";

type Filter = "all" | "unread" | "sales" | "system";

export default function NotificationsPage() {
  const {
    notifications,
    unreadCount,
    isLoading,
    markAsRead,
    markAllAsRead,
    deleteNotification,
    deleteAll,
  } = useNotifications();
  const [filter, setFilter] = useState<Filter>("all");

  const filtered = notifications.filter((n) => {
    if (filter === "unread") return !n.read;
    if (filter === "sales") return n.type === "sale";
    if (filter === "system") return n.type === "system";
    return true;
  });

  return (
    <div className="space-y-6 animate-fade-in max-w-3xl">
      <PageHeader title="Centro de Notificaciones" description="Gestiona tus alertas de ventas y del sistema.">
        <Button variant="outline" onClick={markAllAsRead} disabled={unreadCount === 0}>
          <IconCheck className="size-4 mr-2" /> Marcar todas como leídas
        </Button>
        <Button variant="outline" onClick={deleteAll} disabled={notifications.length === 0}>
          <IconTrash className="size-4 mr-2" /> Vaciar bandeja
        </Button>
      </PageHeader>

      {/* Pestañas de filtro */}
      <div className="flex flex-wrap gap-1.5">
        {[
          { id: "all", label: "Todas", count: notifications.length },
          { id: "unread", label: "No leídas", count: unreadCount },
          { id: "sales", label: "Ventas", count: notifications.filter((n) => n.type === "sale").length },
          { id: "system", label: "Sistema", count: notifications.filter((n) => n.type === "system").length },
        ].map((tab) => (
          <button
            key={tab.id}
            onClick={() => setFilter(tab.id as Filter)}
            className={`rounded-full px-3 py-1.5 text-sm font-medium transition-colors ${
              filter === tab.id
                ? "bg-foreground text-background"
                : "bg-muted text-muted-foreground hover:bg-muted/70"
            }`}
          >
            {tab.label}
            {tab.count > 0 && (
              <span className="ml-1.5 text-xs opacity-70">{tab.count}</span>
            )}
          </button>
        ))}
      </div>

      {/* Lista de notificaciones */}
      <Card className="rounded-2xl border-border shadow-sm">
        <CardContent className="p-0">
          {isLoading ? (
            <div className="p-4 space-y-3">
              {[1, 2, 3].map((i) => <Skeleton key={i} className="h-16 w-full" />)}
            </div>
          ) : filtered.length === 0 ? (
            <div className="p-12 text-center">
              <div className="mx-auto w-12 h-12 rounded-full bg-muted flex items-center justify-center mb-3">
                <IconBell className="size-6 text-muted-foreground" />
              </div>
              <p className="text-sm text-muted-foreground">
                {filter === "all" ? "No tienes notificaciones todavía." : "No hay notificaciones en esta categoría."}
              </p>
            </div>
          ) : (
            <div className="divide-y divide-border">
              {filtered.map((n) => (
                <div
                  key={n.id}
                  className={`flex items-start gap-3 p-4 hover:bg-accent/30 transition-colors ${!n.read ? "bg-accent/10" : ""}`}
                >
                  {!n.read ? (
                    <span className="mt-2 size-2 rounded-full bg-destructive shrink-0" />
                  ) : (
                    <span className="mt-2 size-2 rounded-full bg-transparent shrink-0" />
                  )}

                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 flex-wrap">
                      <Badge
                        variant="outline"
                        className={
                          n.type === "sale"
                            ? "bg-[#F0F5E6] text-[#95B54D] border-[#A7C957]/20"
                            : "bg-neutral-100 text-neutral-600 border-neutral-200"
                        }
                      >
                        {n.type === "sale" ? "Venta" : "Sistema"}
                      </Badge>
                      <span className="text-xs text-muted-foreground">{timeAgo(n.createdAt)}</span>
                    </div>
                    <p className="text-sm font-medium mt-1">{n.title}</p>
                    <p className="text-sm text-muted-foreground">{n.description}</p>
                    <Link
                      href={n.href}
                      className="inline-flex items-center gap-1 text-xs text-foreground font-medium mt-1 hover:underline"
                    >
                      Ver detalle <IconArrowRight className="size-3" />
                    </Link>
                  </div>

                  <div className="flex items-center gap-1 shrink-0">
                    {!n.read && (
                      <Button
                        variant="ghost"
                        size="icon"
                        className="size-8"
                        onClick={() => markAsRead(n.id)}
                        aria-label="Marcar como leída"
                      >
                        <IconCheck className="size-4" />
                      </Button>
                    )}
                    <Button
                      variant="ghost"
                      size="icon"
                      className="size-8 text-muted-foreground hover:text-destructive"
                      onClick={() => deleteNotification(n.id)}
                      aria-label="Eliminar"
                    >
                      <IconTrash className="size-4" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Preferencias de correo (mantenidas del diseño original) */}
      <Card className="rounded-2xl border-border shadow-sm">
        <CardHeader><CardTitle>Preferencias de correo</CardTitle></CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium">Nuevos pedidos</p>
              <p className="text-xs text-muted-foreground">Recibe un correo cada vez que entre una nueva venta.</p>
            </div>
            <input type="checkbox" defaultChecked className="w-4 h-4 accent-foreground" />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}