// core/guards/role.guard.ts
import { CanActivateFn } from '@angular/router';
import { inject } from '@angular/core';
import { AuthService } from '../core/services/auth.service';
import { Router } from '@angular/router';

export const roleGuard: CanActivateFn = (route, state) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  const expectedRoles = route.data['roles'];
  const userRole = authService.getCurrentUserRole();

  if (!expectedRoles.includes(userRole)) {
    return router.parseUrl('/');
  }
  return true;
};
