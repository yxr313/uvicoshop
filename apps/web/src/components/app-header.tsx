"use client";

import * as React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { IconBell, IconHelpCircle, IconTrash } from "@tabler/icons-react";

import { useNotifications, timeAgo } from "@/hooks/use-notifications";
import { SidebarTrigger } from "@/components/ui/sidebar";
import { Button } from "@/components/ui/button";
import {
  Breadcrumb, BreadcrumbItem, BreadcrumbList, BreadcrumbPage, BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import {
  DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

const CRUMB_LABELS: Record<string, string> = {
  hr: "HR", crm: "CRM", api: "API", otp: "OTP", faq: "FAQ",
};

function useBreadcrumb() {
  const pathname = usePathname();
  const segments = pathname.split("/").filter(Boolean);
  return segments.map(
    (s) => CRUMB_LABELS[s] ?? s.replace(/-/g, " ").replace(/\b\w/g, (c) => c.toUpperCase())
  );
}

export function AppHeader() {
  const crumbs = useBreadcrumb();

  return (
    <header className="sticky top-0 z-30 flex h-16 shrink-0 items-center gap-2 border-b bg-background/80 px-4 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <SidebarTrigger className="-ml-1" />
      <div className="mr-1 h-5 w-px shrink-0 self-center bg-border" />
      <Breadcrumb className="hidden sm:block">
        <BreadcrumbList>
          {crumbs.map((c, i) => (
            <React.Fragment key={i}>
              <BreadcrumbItem>
                <BreadcrumbPage className={i < crumbs.length - 1 ? "text-muted-foreground" : ""}>
                  {c}
                </BreadcrumbPage>
              </BreadcrumbItem>
              {i < crumbs.length - 1 && <BreadcrumbSeparator />}
            </React.Fragment>
          ))}
        </BreadcrumbList>
      </Breadcrumb>

      <div className="ml-auto flex items-center gap-2">
        <NotificationsMenu />
        <Button
          variant="ghost"
          size="icon"
          className="rounded-full text-muted-foreground hover:text-foreground"
          aria-label="Ayuda y soporte"
        >
          <IconHelpCircle className="size-5" />
        </Button>
      </div>
    </header>
  );
}

function NotificationsMenu() {
  const { notifications, unreadCount, isLoading, markAsRead, markAllAsRead, deleteNotification, deleteAll } = useNotifications();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          className="relative rounded-full text-muted-foreground hover:text-foreground"
          aria-label="Notificaciones"
        >
          <IconBell className="size-5" />
          {/* Solo mostrar el contador cuando ya cargaron los datos */}
          {!isLoading && unreadCount > 0 && (
            <span className="absolute right-1 top-1 flex size-4 items-center justify-center rounded-full bg-destructive text-[9px] font-semibold text-white ring-2 ring-background">
              {unreadCount > 9 ? "9+" : unreadCount}
            </span>
          )}
        </Button>
      </DropdownMenuTrigger>

      <DropdownMenuContent align="end" className="w-[calc(100vw-2rem)] sm:w-[380px]">
        <div className="flex items-center justify-between px-2 py-1.5">
          <DropdownMenuLabel className="p-0">Notificaciones</DropdownMenuLabel>
          <div className="flex items-center gap-1">
            {!isLoading && unreadCount > 0 && (
              <Button variant="ghost" size="sm" className="h-7 text-xs" onClick={markAllAsRead}>
                Marcar leídas
              </Button>
            )}
            {!isLoading && notifications.length > 0 && (
              <Button variant="ghost" size="icon" className="size-7" onClick={deleteAll} aria-label="Eliminar todas">
                <IconTrash className="size-3.5" />
              </Button>
            )}
          </div>
        </div>
        <DropdownMenuSeparator />

        {isLoading ? (
          <div className="p-6 text-center text-sm text-muted-foreground">
            Cargando...
          </div>
        ) : notifications.length === 0 ? (
          <div className="p-6 text-center text-sm text-muted-foreground">
            Estás al día. No tienes notificaciones.
          </div>
        ) : (
          <div className="max-h-[400px] overflow-y-auto">
            {notifications.slice(0, 6).map((n) => (
              <div
                key={n.id}
                className={`group relative flex items-start gap-3 px-3 py-2.5 hover:bg-accent/50 transition-colors cursor-pointer ${!n.read ? "bg-accent/20" : ""}`}
                onClick={() => markAsRead(n.id)}
              >
                {!n.read && (
                  <span className="absolute left-1 top-1/2 -translate-y-1/2 size-1.5 rounded-full bg-destructive" />
                )}
                <div className="grid min-w-0 flex-1 gap-0.5 pl-2">
                  <p className="line-clamp-2 text-sm leading-snug font-medium">{n.title}</p>
                  <p className="text-xs text-muted-foreground">{n.description}</p>
                  <span className="text-[10px] text-muted-foreground/70 mt-0.5">{timeAgo(n.createdAt)}</span>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  className="size-6 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
                  onClick={(e) => {
                    e.stopPropagation();
                    deleteNotification(n.id);
                  }}
                  aria-label="Eliminar notificación"
                >
                  <IconTrash className="size-3.5" />
                </Button>
              </div>
            ))}
          </div>
        )}

        <DropdownMenuSeparator />
        <DropdownMenuItem asChild className="justify-center text-sm font-medium cursor-pointer">
          <Link href="/dashboard/account/notifications">Ver centro de notificaciones</Link>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}