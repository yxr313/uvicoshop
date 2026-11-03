export type SupplierCategory = "Ropa" | "Calzado" | "Tecnología" | "Hogar" | "Belleza" | "Deportes" | "Juguetes" | "Alimentos" | "Accesorios" | "Papelería";

export type SupplierState = "Distrito Capital" | "Miranda" | "Carabobo" | "Zulia" | "Lara" | "Anzoátegui" | "Aragua" | "Táchira" | "Bolívar" | "Nueva Esparta" | "Falcón" | "Yaracuy" | "Apure";

export interface SupplierProduct {
  id: string;
  name: string;
  description: string;
  priceCents: number;
  imageUrl?: string;
  rating: number; // 1 a 5
  monthlyRotation: number; // Unidades vendidas al mes
  marketAdoption: number; // Porcentaje (ej. 85%)
}

export interface SupplierDocument {
  id: string;
  title: string;
  type: "RIF" | "Registro Mercantil" | "Políticas" | "Contrato";
  url: string; // URL simulada del PDF
}

export interface Supplier {
  id: string;
  name: string;
  description: string;
  category: SupplierCategory;
  state: SupplierState;
  address: string;
  rif: string;
  isVerified: boolean;
  trustScore: number; // 1.0 a 5.0
  slaHours: number;
  logoInitials: string;
  
  // NUEVOS CAMPOS (Red Social B2B)
  bannerUrl?: string;
  avatarUrl?: string;
  connectionsCount: number; // "Seguidores"
  productsCount: number; // "Publicaciones"
  
  // Reputación
  completedOrders: number;
  responseTimeHours: number;
  returnRate: number; // Porcentaje de devoluciones
  
  // Catálogo y Documentos
  products: SupplierProduct[];
  documents: SupplierDocument[];
}