import {
  IconLayoutGrid,
  IconLayoutDashboard,
  IconChartBar,
  IconChartPie,
  IconBuildingStore,
  IconBox,
  IconShoppingCart,
  IconUsers,
  IconBriefcase,
  IconUsersGroup,
  IconTruckDelivery,
  IconLayout,
  IconUserCircle,
  IconCreditCard,
  IconSettings,
  IconShieldLock,
  IconBell,
  IconSparkles,
  IconRocket,
  IconWorld,
  IconBuildingFactory, // NUEVO: Icono para Proveedores
  IconPalette,         // NUEVO: Icono para POD
  type Icon,
} from "@tabler/icons-react";

// ─── Workspaces ──────────────────────────────────────────────────
export type Workspace = {
  name: string;
  plan: string;
  icon: Icon;
};

export const workspaces: Workspace[] = [
  { name: "Mi Tienda", plan: "Fundador", icon: IconSparkles },
  { name: "Tienda Norte", plan: "Pro", icon: IconRocket },
  { name: "Tienda Sur", plan: "Free", icon: IconWorld },
];

// ─── Tipos de navegación ─────────────────────────────────────────
export type NavChild = { title: string; url: string };

export type NavItem = {
  title: string;
  url?: string;
  icon?: Icon;
  badge?: string;
  items?: NavChild[];
};

export type NavGroup = {
  label: string;
  icon: Icon;
  items: NavItem[];
};

// ─── Navegación ──────────────────────────────────────────────────
export const navGroups: NavGroup[] = [
  {
    label: "Panel",
    icon: IconLayoutGrid,
    items: [
      { title: "Overview", url: "/dashboard", icon: IconLayoutDashboard },
      { title: "Analytics", url: "/dashboard/analytics", icon: IconChartBar },
    ],
  },
  {
    label: "Mi Tienda",
    icon: IconBuildingStore,
    items: [
      { title: "Resumen", url: "/dashboard/ecommerce", icon: IconChartPie },
      { title: "Productos", url: "/dashboard/ecommerce/products", icon: IconBox },
      { title: "Pedidos", url: "/dashboard/ecommerce/orders", icon: IconShoppingCart },
      { title: "Clientes", url: "/dashboard/ecommerce/customers", icon: IconUsers },
    ],
  },
  {
    label: "Proveedores",
    icon: IconBuildingFactory,
    items: [
      {
        title: "Dropshipping",
        icon: IconTruckDelivery,
        items: [
          { title: "Directorio", url: "/dashboard/suppliers/dropshipping" },
          { title: "Mis Conexiones", url: "/dashboard/suppliers/dropshipping/connections" },
          { title: "Pedidos de Prueba", url: "/dashboard/suppliers/dropshipping/samples" },
        ],
      },
      {
        title: "Print on Demand",
        icon: IconPalette,
        items: [
          { title: "Estudio de Diseño", url: "/dashboard/suppliers/pod" },
          { title: "Mis Diseños", url: "/dashboard/suppliers/pod/designs" },
          { title: "Órdenes POD", url: "/dashboard/suppliers/pod/orders" },
        ],
      },
    ],
  },
  {
    label: "Operaciones",
    icon: IconBriefcase,
    items: [
      { title: "CRM", url: "/dashboard/crm", icon: IconUsersGroup },
      { title: "Logística", url: "/dashboard/logistics", icon: IconTruckDelivery },
      { title: "Editor Visual", url: "/dashboard/editor", icon: IconLayout },
    ],
  },
  {
    label: "Cuenta",
    icon: IconUserCircle,
    items: [
      { title: "Suscripción", url: "/dashboard/pricing", icon: IconCreditCard },
      { title: "Ajustes", url: "/dashboard/settings", icon: IconSettings },
      { title: "Perfil", url: "/dashboard/account/profile", icon: IconUserCircle },
      { title: "Seguridad", url: "/dashboard/account/security", icon: IconShieldLock },
      { title: "Notificaciones", url: "/dashboard/account/notifications", icon: IconBell },
    ],
  },
];