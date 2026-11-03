import { TooltipProvider } from "@/components/ui/tooltip";
import { SidebarProvider, SidebarInset } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/app-sidebar";
import { AppHeader } from "@/components/app-header";
import { QueryProvider } from "./query-provider"; // <-- Importamos el motor
import "@puckeditor/core/puck.css";

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  return (
    <QueryProvider> {/* <-- Envolvemos todo */}
      <TooltipProvider delayDuration={200}>
        <SidebarProvider defaultOpen>
          <AppSidebar />
          <SidebarInset>
            <AppHeader />
            <main className="flex-1 p-4 md:p-6">{children}</main>
          </SidebarInset>
        </SidebarProvider>
      </TooltipProvider>
    </QueryProvider>
  );
}