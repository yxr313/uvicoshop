import { useMemo, useEffect, useState } from "react";
import { useDashboardData } from "./use-dashboard";

export type NotificationType = "sale" | "system";

export interface Notification {
  id: string;
  type: NotificationType;
  title: string;
  description: string;
  amountCents?: number;
  createdAt: string;
  read: boolean;
  deleted: boolean;
  href: string;
}

const STORAGE_KEY = "uvico-notifications-state";

interface NotificationState {
  read: string[];
  deleted: string[];
}

function loadState(): NotificationState {
  if (typeof window === "undefined") return { read: [], deleted: [] };
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) return JSON.parse(raw);
  } catch {
    // Ignorar errores de parseo
  }
  return { read: [], deleted: [] };
}

function saveState(state: NotificationState) {
  if (typeof window === "undefined") return;
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
  } catch {
    // Ignorar errores de almacenamiento
  }
}

export function useNotifications() {
  const { recentOrders, isLoading: isLoadingOrders } = useDashboardData();
  const [state, setState] = useState<NotificationState>({ read: [], deleted: [] });
  const [hydrated, setHydrated] = useState(false);

  // Hidratar desde localStorage al montar
  useEffect(() => {
    setState(loadState());
    setHydrated(true);
  }, []);

  const notifications: Notification[] = useMemo(() => {
    // 1. Notificaciones derivadas de ventas reales
    const fromOrders: Notification[] = recentOrders.map((o) => {
      const notificationId = `sale-${o.id}`;
      return {
        id: notificationId,
        type: "sale" as const,
        title: `Nueva venta de ${o.customer_name || "Cliente"}`,
        description: `Pedido por $${(o.total_cents / 100).toFixed(2)}`,
        amountCents: o.total_cents,
        createdAt: o.created_at,
        read: state.read.includes(notificationId),
        deleted: state.deleted.includes(notificationId),
        href: `/dashboard/ecommerce/orders`,
      };
    });

    // 2. Notificaciones del sistema
    // Usamos una fecha fija o calculada de forma idéntica en servidor y cliente para evitar hydration mismatch
    const systemDate = "2026-10-31T12:00:00.000Z"; 
    const system: Notification[] = [
      {
        id: "system-welcome",
        type: "system" as const,
        title: "Bienvenido a Uvico",
        description: "Tu tienda está lista. Configura tus métodos de pago para empezar a vender.",
        createdAt: systemDate,
        read: state.read.includes("system-welcome"),
        deleted: state.deleted.includes("system-welcome"),
        href: "/dashboard/settings",
      },
    ];

    return [...fromOrders, ...system]
      .filter((n) => !n.deleted)
      .sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
  }, [recentOrders, state]);

  const unreadCount = notifications.filter((n) => !n.read).length;

  const markAsRead = (id: string) => {
    setState((prev) => {
      if (prev.read.includes(id)) return prev;
      const next = { ...prev, read: [...prev.read, id] };
      saveState(next);
      return next;
    });
  };

  const markAllAsRead = () => {
    setState((prev) => {
      const allIds = notifications.map((n) => n.id);
      const next = { ...prev, read: Array.from(new Set([...prev.read, ...allIds])) };
      saveState(next);
      return next;
    });
  };

  const deleteNotification = (id: string) => {
    setState((prev) => {
      const next = { ...prev, deleted: [...prev.deleted, id] };
      saveState(next);
      return next;
    });
  };

  const deleteAll = () => {
    setState((prev) => {
      const allIds = notifications.map((n) => n.id);
      const next = { ...prev, deleted: Array.from(new Set([...prev.deleted, ...allIds])) };
      saveState(next);
      return next;
    });
  };

  return {
    notifications,
    unreadCount,
    isLoading: isLoadingOrders || !hydrated,
    markRead: markAsRead, // Corregido nombre de exportación
    markAsRead,
    markAllAsRead,
    deleteNotification,
    deleteAll,
  };
}

// Helper para mostrar tiempo relativo
export function timeAgo(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 1) return "Hace un momento";
  if (mins < 60) return `Hace ${mins} min`;
  const hours = Math.floor(mins / 60);
  if (hours < 24) return `Hace ${hours}h`;
  const days = Math.floor(hours / 24);
  if (days < 7) return `Hace ${days}d`;
  return new Date(iso).toLocaleDateString("es-VE", { day: "numeric", month: "short" });
}