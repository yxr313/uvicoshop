import { 
  LayoutDashboard, BarChart3, ShoppingCart, Users, 
  Settings, ShieldCheck, Bell, Package, ListOrdered, UserPlus,
  FileEdit, CreditCard, UserCircle
} from "lucide-react";

export const sidebarData = [
  {
    title: "Dashboards",
    items: [
      { title: "Overview", url: "/dashboard", icon: LayoutDashboard },
      { title: "Analytics", url: "/dashboard/analytics", icon: BarChart3 },
      { title: "E-Commerce", url: "/dashboard/ecommerce", icon: ShoppingCart },
      { title: "CRM", url: "/dashboard/crm", icon: Users },
      { title: "Editor Visual", url: "/dashboard/editor", icon: FileEdit }, // Movido aquí
    ],
  },
  {
    title: "E-commerce",
    items: [
      {
        title: "Tienda",
        icon: ShoppingCart,
        items: [
          { title: "Productos", url: "/dashboard/ecommerce/products", icon: Package },
          { title: "Pedidos", url: "/dashboard/ecommerce/orders", icon: ListOrdered },
          { title: "Clientes", url: "/dashboard/ecommerce/customers", icon: UserPlus },
        ],
      },
    ],
  },
  {
    title: "Configuración",
    items: [
      {
        title: "Account",
        icon: UserCircle,
        items: [
          { title: "Perfil", url: "/dashboard/account/profile", icon: UserCircle },
          { title: "Seguridad", url: "/dashboard/account/security", icon: ShieldCheck },
          { title: "Notificaciones", url: "/dashboard/account/notifications", icon: Bell },
        ],
      },
      { title: "Pricing", url: "/dashboard/pricing", icon: CreditCard },
      { title: "Ajustes", url: "/dashboard/settings", icon: Settings },
    ],
  },
];
