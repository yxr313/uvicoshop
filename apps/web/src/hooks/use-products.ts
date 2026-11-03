import { useQuery } from "@tanstack/react-query";
import { fetchFromApi } from "@/lib/api-client";

export interface Product {
  id: string;
  tenant_id: string;
  provider_id?: string;
  sale_mode: "own_inventory" | "dropshipping_verified";
  sku: string;
  slug: string;
  name: string;
  description: string;
  price_cents: number;
  currency: string;
  status: "draft" | "active" | "archived";
  stock: number;
  created_at: string;
  updated_at: string;
  version: number;
}

interface ProductsResponse {
  products: Product[];
}

export function useProducts(filters: { search?: string; mode?: string; status?: string } = {}) {
  return useQuery({
    queryKey: ["products", filters],
    queryFn: async () => {
      const params = new URLSearchParams();
      if (filters.search) params.set("search", filters.search);
      if (filters.mode && filters.mode !== "all") params.set("sale_mode", filters.mode);
      if (filters.status && filters.status !== "all") params.set("status", filters.status);

      const queryString = params.toString() ? `?${params.toString()}` : "";
      const res = await fetchFromApi<ProductsResponse>(`/products${queryString}`);
      return res.products || [];
    },
  });
}