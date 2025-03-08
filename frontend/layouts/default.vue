<template>
  <v-layout class="rounded rounded-md">
    <v-app-bar
      :elevation="1"
      title="Orion DNS"
    />

    <v-navigation-drawer
      expand-on-hover
      rail
    >
      <v-list>
        <v-list-item
          prepend-icon="mdi-view-dashboard"
          subtitle="View your data"
          title="Dashboard"
          to="/dashboard"
        />
      </v-list>
    </v-navigation-drawer>

    <v-main>
      <slot />
    </v-main>
  </v-layout>
</template>

<style>
.v-application__wrap {
  min-height: 0vh !important;
}
</style>
