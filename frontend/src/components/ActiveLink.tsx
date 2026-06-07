"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";

interface ActiveLinkProps extends React.ComponentPropsWithoutRef<typeof Link> {
  activeClassName?: string;
  end?: boolean;
}

export function ActiveLink({
  className,
  activeClassName,
  end = false,
  href,
  ...props
}: ActiveLinkProps) {
  const pathname = usePathname();
  const hrefStr = typeof href === "string" ? href : href.toString();
  const isActive = end ? pathname === hrefStr : hrefStr === "/" ? pathname === "/" : pathname.startsWith(hrefStr);

  return (
    <Link
      href={href}
      className={cn(className, isActive && activeClassName)}
      {...props}
    />
  );
}