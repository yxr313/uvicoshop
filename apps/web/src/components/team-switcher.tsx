"use client";

import * as React from "react";
import { IconSelector, IconCheck, IconBuildingStore } from "@tabler/icons-react";

import { workspaces, type Workspace } from "@/config/nav";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar";

const FALLBACK_WORKSPACE: Workspace = {
  name: "Mi Tienda",
  plan: "Fundador",
  icon: IconBuildingStore,
};

export function TeamSwitcher() {
  const { isMobile } = useSidebar();
  const [active, setActive] = React.useState<Workspace>(
    workspaces[0] ?? FALLBACK_WORKSPACE
  );

  const Icon = active.icon;

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <div className="flex aspect-square size-9 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                <Icon className="size-5" />
              </div>
              <div className="grid flex-1 text-left leading-tight">
                <span className="truncate font-semibold tracking-tight">
                  {active.name}
                </span>
                <span className="truncate text-xs text-muted-foreground">
                  {active.plan}
                </span>
              </div>
              <IconSelector className="ml-auto size-4 text-muted-foreground" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-(--radix-dropdown-menu-trigger-width) min-w-60 rounded-lg"
            align="start"
            side={isMobile ? "bottom" : "right"}
            sideOffset={4}
          >
            <DropdownMenuLabel className="text-xs text-muted-foreground">
              Tiendas
            </DropdownMenuLabel>
            {workspaces.map((ws) => {
              const WsIcon = ws.icon;
              return (
                <DropdownMenuItem
                  key={ws.name}
                  onClick={() => setActive(ws)}
                  className="gap-2 p-2"
                >
                  <div className="flex size-6 items-center justify-center rounded-md border">
                    <WsIcon className="size-3.5 shrink-0" />
                  </div>
                  <div className="grid flex-1 leading-tight">
                    <span className="truncate text-sm font-medium">{ws.name}</span>
                    <span className="truncate text-xs text-muted-foreground">
                      {ws.plan}
                    </span>
                  </div>
                  {active.name === ws.name && (
                    <IconCheck className="size-4 shrink-0" />
                  )}
                </DropdownMenuItem>
              );
            })}
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}