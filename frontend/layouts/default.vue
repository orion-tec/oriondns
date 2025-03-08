<template>
  <v-layout class="rounded rounded-md">
    <v-app-bar
      title="Orion DNS"
      :elevation="1"
    />

    <v-navigation-drawer
      expand-on-hover
      rail
    >
      <v-list>
        <v-list-item
          prepend-icon="mdi-view-dashboard"
          to="/dashboard"
          title="Dashboard"
          subtitle="View your data"
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
