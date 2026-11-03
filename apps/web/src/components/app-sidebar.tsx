"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { IconChevronRight } from "@tabler/icons-react";

import { cn } from "@/lib/utils";
import { navGroups, type NavItem } from "@/config/nav";
import { NavUser } from "@/components/nav-user";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuBadge,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarRail,
  useSidebar,
} from "@/components/ui/sidebar";

export function AppSidebar() {
  const pathname = usePathname();

  return (
    <Sidebar collapsible="icon" variant="sidebar" className="border-r border-border bg-white">
      
      <SidebarHeader className="h-16 flex items-center justify-center border-b border-border px-6 group-data-[collapsible=icon]:px-0">
        <Link href="/" className="font-display font-bold text-xl tracking-tight text-foreground group-data-[collapsible=icon]:hidden">
          UVICO<span className="text-accent">.</span>
        </Link>
        <Link href="/" className="font-display font-bold text-xl tracking-tight text-foreground hidden group-data-[collapsible=icon]:block">
          U<span className="text-accent">.</span>
        </Link>
      </SidebarHeader>

      <SidebarContent className="py-6 px-3 no-scrollbar">
        {navGroups.map((group) => (
          <SidebarGroup key={group.label} className="p-0 mb-6">
            <SidebarGroupLabel className="h-auto px-3 mb-2 text-[11px] font-semibold uppercase tracking-wider text-neutral-400 group-data-[collapsible=icon]:hidden flex items-center gap-2">
              <group.icon className="size-3.5 shrink-0" />
              {group.label}
            </SidebarGroupLabel>
            <SidebarMenu className="gap-1">
              {group.items.map((item) => (
                <NavEntry key={item.title} item={item} pathname={pathname} />
              ))}
            </SidebarMenu>
          </SidebarGroup>
        ))}
      </SidebarContent>

      <SidebarFooter className="p-4 border-t border-border">
        <NavUser />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}

function badgeClass(badge: string) {
  return /^\d+$/.test(badge)
    ? "bg-neutral-100 text-neutral-600"
    : "bg-accent px-1.5 text-[10px] font-semibold text-accent-foreground";
}

function NavEntry({ item, pathname }: { item: NavItem; pathname: string }) {
  const { setOpenMobile } = useSidebar();
  const closeMobile = () => setOpenMobile(false);

  if (item.items?.length) {
    const childActive = item.items.some((c) => pathname === c.url);
    return (
      <Collapsible
        asChild
        defaultOpen={childActive}
        className="group/collapsible"
      >
        <SidebarMenuItem>
          <CollapsibleTrigger asChild>
            <SidebarMenuButton
              tooltip={item.title}
              isActive={childActive}
              className={cn(
                "h-9 gap-3 px-3 py-2 text-sm font-medium transition-colors rounded-md",
                childActive ? "bg-neutral-100 text-foreground" : "text-muted-foreground hover:bg-neutral-100 hover:text-foreground"
              )}
            >
              {item.icon && <item.icon className="size-[18px] shrink-0" />}
              <span>{item.title}</span>
              {item.badge && (
                <SidebarMenuBadge
                  className={cn(
                    "group-data-[state=open]/collapsible:hidden",
                    badgeClass(item.badge)
                  )}
                >
                  {item.badge}
                </SidebarMenuBadge>
              )}
              <IconChevronRight className="ml-auto size-4 transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
            </SidebarMenuButton>
          </CollapsibleTrigger>
          <CollapsibleContent className="overflow-hidden [--tw-animation-duration:260ms] data-[state=closed]:animate-collapsible-up data-[state=open]:animate-collapsible-down">
            <SidebarMenuSub className="border-l border-neutral-200 ml-5 pl-4">
              {item.items.map((sub) => (
                <SidebarMenuSubItem key={sub.title}>
                  <SidebarMenuSubButton
                    asChild
                    isActive={pathname === sub.url}
                    className={cn(
                      "h-8 px-2 text-sm rounded-md transition-colors",
                      pathname === sub.url ? "bg-neutral-100 text-foreground font-medium" : "text-muted-foreground hover:bg-neutral-100 hover:text-foreground"
                    )}
                  >
                    <Link href={sub.url} onClick={closeMobile}>
                      {sub.title}
                    </Link>
                  </SidebarMenuSubButton>
                </SidebarMenuSubItem>
              ))}
            </SidebarMenuSub>
          </CollapsibleContent>
        </SidebarMenuItem>
      </Collapsible>
    );
  }

  const active =
    pathname === item.url ||
    (item.url !== "/dashboard" && !!item.url && pathname.startsWith(item.url));
    
  return (
    <SidebarMenuItem>
      <SidebarMenuButton
        asChild
        isActive={active}
        tooltip={item.title}
        className={cn(
          "h-9 gap-3 px-3 py-2 text-sm font-medium transition-colors rounded-md",
          active ? "bg-neutral-100 text-foreground" : "text-muted-foreground hover:bg-neutral-100 hover:text-foreground"
        )}
      >
        <Link href={item.url ?? "#"} onClick={closeMobile}>
          {item.icon && <item.icon className="size-[18px] shrink-0" />}
          <span>{item.title}</span>
        </Link>
      </SidebarMenuButton>
      {item.badge && (
        <SidebarMenuBadge className={badgeClass(item.badge)}>
          {item.badge}
        </SidebarMenuBadge>
      )}
    </SidebarMenuItem>
  );
}