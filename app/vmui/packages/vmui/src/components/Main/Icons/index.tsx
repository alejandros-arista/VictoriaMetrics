import React from "react";
import { getCssVariable } from "../../../utils/theme";

export const LogoIcon = () => (
  <svg
    viewBox="0 0 74 24"
    fill="currentColor"
  >
    <path
      d="M6.12 10.48c.36.28.8.43 1.26.43h.05c.48 0 .96-.19 1.25-.44 1.5-1.28 5.88-5.29 5.88-5.29C15.73 4.1 12.46 3.01 7.43 3h-.06C2.33 3-.93 4.1.24 5.18c0 0 4.37 4 5.88 5.3Zm2.56 2.16c-.36.28-.8.44-1.26.45h-.04c-.46 0-.9-.17-1.26-.45-1.04-.88-4.74-4.22-6.12-5.5v1.94c0 .21.08.5.22.63l.07.06c1.05.96 4.55 4.16 5.83 5.25.36.28.8.43 1.26.44h.04c.49-.02.96-.2 1.26-.44 1.3-1.11 4.94-4.45 5.88-5.31.15-.14.23-.42.23-.63V7.15a454.94 454.94 0 0 1-6.11 5.5Zm-1.26 4.99c.46 0 .9-.16 1.26-.44a454.4 454.4 0 0 0 6.1-5.5v1.94c0 .2-.07.48-.22.62-.94.87-4.57 4.2-5.88 5.3-.3.26-.77.44-1.26.45h-.04c-.46 0-.9-.16-1.26-.44-1.2-1.02-4.38-3.92-5.62-5.06l-.28-.25c-.14-.14-.22-.42-.22-.62v-1.94c1.38 1.26 5.08 4.6 6.12 5.5.36.28.8.43 1.26.44h.04ZM35 5l-5.84 14.46h-2.43L20.89 5h2.16a.9.9 0 0 1 .9.61l3.41 8.82a18.8 18.8 0 0 1 .62 2.02 19.44 19.44 0 0 1 .57-2.02l3.39-8.82c.05-.15.16-.3.31-.42a.9.9 0 0 1 .58-.19H35Zm17.18 0v14.46H49.8v-9.34c0-.37.02-.78.06-1.21l-4.37 8.21c-.21.4-.53.59-.95.59h-.38c-.43 0-.75-.2-.95-.59L38.8 8.88a22.96 22.96 0 0 1 .07 1.24v9.34H36.5V5h2.03l.3.01c.1 0 .17.02.24.05.07.03.13.07.19.13a1 1 0 0 1 .17.24l4.33 8.03a16.97 16.97 0 0 1 .6 1.36 14.34 14.34 0 0 1 .6-1.38l4.28-8.01c.05-.1.1-.18.17-.24.06-.06.12-.1.19-.13a.9.9 0 0 1 .24-.05l.3-.01h2.04Zm8.88 13.73a4.5 4.5 0 0 0 1.82-.35 3.96 3.96 0 0 0 2.22-2.47c.2-.57.3-1.19.3-1.85V5.31h1.02v8.75c0 .78-.12 1.51-.37 2.19a4.88 4.88 0 0 1-2.76 2.95c-.66.29-1.4.43-2.23.43-.82 0-1.57-.14-2.24-.43a5.01 5.01 0 0 1-2.75-2.95 6.37 6.37 0 0 1-.37-2.19V5.31h1.03v8.74c0 .66.1 1.28.3 1.85a3.98 3.98 0 0 0 2.21 2.47c.53.24 1.14.36 1.82.36Zm10.38.73h-1.03V5.31h1.03v14.15Z"
    />
  </svg>
);

export const LogoLogsIcon = () => (
  <svg
    viewBox="0 0 85 38"
    fill="currentColor"
  >
    <path
      d="M11.12 10.48c.36.28.8.43 1.26.43h.05c.48 0 .96-.19 1.25-.44 1.5-1.28 5.88-5.29 5.88-5.29 1.17-1.09-2.1-2.17-7.13-2.18h-.06c-5.04 0-8.3 1.1-7.13 2.18 0 0 4.37 4 5.88 5.3Zm2.56 2.16c-.36.28-.8.44-1.26.45h-.04c-.46 0-.9-.17-1.26-.45-1.04-.88-4.74-4.22-6.12-5.5v1.94c0 .21.08.5.22.63l.07.06c1.05.96 4.55 4.16 5.83 5.25.36.28.8.43 1.26.44h.04c.49-.02.96-.2 1.26-.44 1.3-1.11 4.94-4.45 5.88-5.31.15-.14.23-.42.23-.63V7.15a455.13 455.13 0 0 1-6.11 5.5Zm-1.26 4.99c.46 0 .9-.16 1.26-.44 2.05-1.82 4.09-3.65 6.1-5.5v1.94c0 .2-.07.48-.22.62-.94.87-4.57 4.2-5.88 5.3-.3.26-.77.44-1.26.45h-.04c-.46 0-.9-.16-1.26-.44-1.2-1.02-4.38-3.92-5.62-5.06l-.28-.25c-.14-.14-.22-.42-.22-.62v-1.94c1.38 1.26 5.08 4.6 6.12 5.5.36.28.8.43 1.26.44h.04ZM40 5l-5.84 14.46h-2.43L25.89 5h2.16a.9.9 0 0 1 .9.61l3.41 8.82a18.8 18.8 0 0 1 .62 2.02 19.44 19.44 0 0 1 .57-2.02l3.39-8.82c.05-.15.16-.3.31-.42a.9.9 0 0 1 .58-.19H40Zm17.18 0v14.46H54.8v-9.34c0-.37.02-.78.06-1.21l-4.37 8.21c-.21.4-.53.59-.95.59h-.38c-.43 0-.75-.2-.95-.59L43.8 8.88a22.96 22.96 0 0 1 .07 1.24v9.34H41.5V5h2.03l.3.01c.1 0 .17.02.24.05.07.03.13.07.19.13a1 1 0 0 1 .17.24l4.33 8.03a16.97 16.97 0 0 1 .6 1.36 14.34 14.34 0 0 1 .6-1.38l4.28-8.01c.05-.1.1-.18.17-.24.06-.06.12-.1.19-.13a.9.9 0 0 1 .24-.05l.3-.01h2.04Zm8.88 13.73a4.5 4.5 0 0 0 1.82-.35 3.96 3.96 0 0 0 2.22-2.47c.2-.57.3-1.19.3-1.85V5.31h1.02v8.75c0 .78-.12 1.51-.37 2.19a4.88 4.88 0 0 1-2.76 2.95c-.66.29-1.4.43-2.23.43-.82 0-1.57-.14-2.24-.43a5.01 5.01 0 0 1-2.75-2.95 6.37 6.37 0 0 1-.37-2.19V5.31h1.03v8.74c0 .66.1 1.28.3 1.85a3.98 3.98 0 0 0 2.21 2.47c.53.24 1.14.36 1.82.36Zm10.38.73h-1.03V5.31h1.03v14.15ZM1.73 36v-5.17l-.67-.07a.6.6 0 0 1-.21-.1.23.23 0 0 1-.08-.18v-.44h.96v-.59c0-.34.05-.65.14-.92a1.79 1.79 0 0 1 1.08-1.11 2.45 2.45 0 0 1 1.62-.02l-.03.53c0 .1-.06.15-.16.16H4c-.18 0-.35.03-.5.08a.95.95 0 0 0-.39.23c-.1.11-.19.25-.25.43-.05.18-.08.4-.08.65v.56h1.75v.78H2.8V36H1.73Zm6.17-6.17c.45 0 .85.07 1.2.22a2.57 2.57 0 0 1 1.5 1.62c.13.38.2.81.2 1.29s-.07.91-.2 1.3a2.57 2.57 0 0 1-1.49 1.61c-.36.14-.76.21-1.2.21-.45 0-.86-.07-1.22-.21a2.57 2.57 0 0 1-1.5-1.62c-.12-.38-.19-.81-.19-1.3 0-.47.07-.9.2-1.28a2.57 2.57 0 0 1 1.5-1.62c.35-.15.76-.22 1.2-.22Zm0 5.42c.6 0 1.05-.2 1.35-.6.3-.4.44-.97.44-1.69s-.15-1.28-.44-1.69c-.3-.4-.75-.6-1.35-.6-.3 0-.57.05-.8.15-.22.1-.4.26-.56.45-.15.2-.26.44-.33.73-.08.28-.11.6-.11.96 0 .72.15 1.29.44 1.69.3.4.76.6 1.36.6Zm5.26-4.11c.2-.42.43-.74.71-.97.28-.24.62-.36 1.03-.36.13 0 .25.02.36.05.12.02.23.07.32.13l-.08.8c-.02.1-.08.15-.18.15l-.24-.04a1.7 1.7 0 0 0-.88.05c-.15.05-.29.14-.4.25-.12.1-.23.24-.32.4-.1.17-.18.35-.26.56V36h-1.07v-6.08h.61c.12 0 .2.02.24.07.05.04.08.12.1.23l.06.92Zm13.73-3.82L23.39 36h-1.46l-3.5-8.68h1.29a.54.54 0 0 1 .54.37l2.04 5.3a11.31 11.31 0 0 1 .37 1.21 11.65 11.65 0 0 1 .35-1.22l2.03-5.29c.03-.1.1-.18.19-.25.1-.08.21-.12.35-.12h1.3Zm2.2 2.52V36H27.6v-6.16h1.49Zm.2-1.79c0 .13-.02.25-.08.36a1 1 0 0 1-.51.5.96.96 0 0 1-.73 0 1.02 1.02 0 0 1-.5-.5.96.96 0 0 1 0-.73.93.93 0 0 1 .86-.58.9.9 0 0 1 .37.08c.12.05.22.11.3.2a.94.94 0 0 1 .3.67Zm5.72 3.1a.68.68 0 0 1-.13.13c-.04.03-.1.05-.18.05a.42.42 0 0 1-.22-.07 3.95 3.95 0 0 0-.62-.31c-.14-.05-.3-.07-.51-.07-.26 0-.5.04-.69.14-.2.1-.36.23-.49.4-.13.18-.22.4-.29.64-.06.25-.1.53-.1.85 0 .33.04.62.1.88.08.25.18.47.32.64.13.18.29.3.48.4.18.09.4.13.63.13a1.6 1.6 0 0 0 .94-.27l.26-.2a.4.4 0 0 1 .25-.09.3.3 0 0 1 .27.14l.43.54a2.76 2.76 0 0 1-1.77.96c-.22.03-.43.05-.65.05a2.57 2.57 0 0 1-1.96-.83c-.25-.28-.45-.6-.6-1-.14-.4-.21-.85-.21-1.35 0-.45.06-.87.2-1.25a2.61 2.61 0 0 1 1.51-1.67c.37-.16.8-.24 1.28-.24.46 0 .86.07 1.2.22.35.15.66.36.94.64l-.4.54Zm3.43 4.95c-.54 0-.95-.15-1.24-.45-.28-.3-.42-.73-.42-1.26v-3.44h-.63a.29.29 0 0 1-.2-.07c-.06-.06-.09-.13-.09-.24v-.59l.99-.16.31-1.68a.33.33 0 0 1 .12-.18.34.34 0 0 1 .21-.07h.77v1.94h1.64v1.05h-1.64v3.34c0 .2.05.34.14.45.1.1.22.16.39.16a.73.73 0 0 0 .39-.1l.12-.07a.2.2 0 0 1 .11-.03c.05 0 .08.01.11.03l.09.1.44.72c-.21.18-.46.32-.74.4-.28.1-.57.15-.87.15Zm5.09-6.35c.46 0 .87.07 1.24.22a2.7 2.7 0 0 1 1.58 1.63c.14.39.22.83.22 1.31 0 .49-.08.93-.22 1.32-.14.4-.35.73-.62 1-.26.28-.58.49-.96.64-.37.15-.78.22-1.24.22a3.4 3.4 0 0 1-1.25-.22 2.71 2.71 0 0 1-1.59-1.64 3.8 3.8 0 0 1-.21-1.32c0-.48.07-.92.21-1.31a2.75 2.75 0 0 1 1.58-1.63c.38-.15.8-.22 1.26-.22Zm0 5.2c.51 0 .89-.17 1.13-.52.25-.34.38-.84.38-1.5a2.6 2.6 0 0 0-.38-1.53c-.24-.34-.62-.52-1.13-.52-.52 0-.9.18-1.16.53-.25.35-.37.85-.37 1.51s.12 1.17.37 1.51c.25.35.64.52 1.16.52Zm5.56-4.04c.2-.37.42-.65.69-.86.26-.21.57-.32.94-.32.28 0 .5.06.68.19l-.1 1.1a.3.3 0 0 1-.09.16.24.24 0 0 1-.15.04 1.8 1.8 0 0 1-.27-.03 2.01 2.01 0 0 0-.34-.03c-.16 0-.3.03-.44.08a1.1 1.1 0 0 0-.34.2c-.1.1-.2.2-.27.33-.08.13-.15.27-.22.44V36H47.7v-6.16h.87c.15 0 .26.03.31.09.06.05.1.15.13.29l.09.7Zm4.62-1.07V36h-1.49v-6.16h1.49Zm.2-1.79c0 .13-.02.25-.07.36a1 1 0 0 1-.51.5.96.96 0 0 1-.74 0 1.02 1.02 0 0 1-.5-.5.96.96 0 0 1 0-.73.93.93 0 0 1 .86-.58.9.9 0 0 1 .38.08c.11.05.21.11.3.2a.94.94 0 0 1 .28.67Zm4.56 5.32a7.8 7.8 0 0 0-1.08.12c-.29.05-.52.12-.7.2a.92.92 0 0 0-.38.3.64.64 0 0 0-.11.36c0 .26.07.45.23.56.15.11.35.17.6.17.3 0 .57-.06.79-.17.22-.1.44-.28.65-.5v-1.04Zm-3.4-2.67c.71-.65 1.57-.97 2.56-.97.36 0 .68.06.97.18a1.99 1.99 0 0 1 1.16 1.24c.1.3.16.61.16.96V36h-.67a.7.7 0 0 1-.33-.06c-.07-.04-.13-.13-.18-.26l-.13-.44c-.16.14-.3.26-.46.37a2.8 2.8 0 0 1-.97.43 2.77 2.77 0 0 1-1.32-.05 1.62 1.62 0 0 1-.57-.31 1.41 1.41 0 0 1-.38-.53 1.85 1.85 0 0 1-.05-1.18c.05-.16.14-.3.25-.45.12-.14.28-.27.46-.4a3 3 0 0 1 .7-.32 9.19 9.19 0 0 1 2.2-.33v-.36c0-.41-.09-.71-.26-.91-.18-.2-.43-.3-.76-.3a1.84 1.84 0 0 0-1.02.28l-.33.18c-.1.06-.2.09-.32.09-.1 0-.2-.03-.27-.08a.72.72 0 0 1-.17-.2l-.26-.47Zm11.49 4.32V36h-4.88v-8.6h1.16v7.62h3.72Zm3.16-5.2c.44 0 .84.08 1.2.23a2.57 2.57 0 0 1 1.49 1.62c.13.38.2.81.2 1.29s-.07.91-.2 1.3a2.57 2.57 0 0 1-1.49 1.61c-.36.14-.76.21-1.2.21-.45 0-.85-.07-1.21-.21a2.57 2.57 0 0 1-1.5-1.62c-.13-.38-.2-.81-.2-1.3 0-.47.07-.9.2-1.28.14-.39.33-.72.59-1 .25-.26.55-.47.9-.62.37-.15.77-.22 1.22-.22Zm0 5.43c.6 0 1.05-.2 1.34-.6.3-.4.45-.97.45-1.69s-.15-1.28-.45-1.69c-.3-.4-.74-.6-1.34-.6-.3 0-.57.05-.8.15-.22.1-.4.26-.56.45-.15.2-.26.44-.34.73-.07.28-.1.6-.1.96 0 .72.14 1.29.44 1.69.3.4.75.6 1.36.6Zm6.33-2.22c.22 0 .4-.03.57-.09.16-.06.3-.14.41-.25.12-.11.2-.24.26-.39.05-.15.08-.31.08-.5 0-.37-.11-.66-.34-.88-.23-.22-.55-.33-.98-.33-.43 0-.76.1-.99.33-.22.22-.34.51-.34.89 0 .18.03.34.09.5a1.1 1.1 0 0 0 .67.63c.16.06.35.09.57.09Zm1.93 3.3a.51.51 0 0 0-.13-.36.84.84 0 0 0-.34-.22 8.57 8.57 0 0 0-1.73-.2 7.5 7.5 0 0 1-.62-.05c-.23.1-.41.23-.56.4a.8.8 0 0 0-.1.92c.07.12.18.22.32.3.14.1.32.16.54.21a3.5 3.5 0 0 0 1.55 0c.23-.05.42-.12.57-.22.16-.1.29-.21.37-.34a.8.8 0 0 0 .13-.44Zm1.08-6.17v.4c0 .13-.08.21-.25.25l-.69.09c.14.26.2.56.2.88a1.86 1.86 0 0 1-1.36 1.82 3.07 3.07 0 0 1-1.72.04c-.12.08-.22.16-.29.25a.44.44 0 0 0-.1.27c0 .15.06.26.17.33.12.08.28.13.47.16a5 5 0 0 0 .66.06 16.56 16.56 0 0 1 1.5.13c.26.05.48.12.67.22.19.1.34.24.46.41.12.18.18.4.18.69 0 .26-.07.5-.2.75s-.31.46-.56.65c-.24.2-.54.34-.9.46a4.57 4.57 0 0 1-2.36.04c-.33-.09-.6-.2-.82-.36a1.56 1.56 0 0 1-.5-.51c-.1-.2-.16-.4-.16-.6 0-.3.1-.56.28-.77.19-.2.45-.37.77-.5a1.15 1.15 0 0 1-.43-.32.88.88 0 0 1-.15-.54c0-.09.01-.18.04-.27.04-.1.08-.2.15-.28a1.55 1.55 0 0 1 .58-.5c-.3-.16-.53-.39-.7-.66-.17-.28-.25-.6-.25-.97 0-.3.05-.57.16-.8.12-.25.28-.46.48-.63.2-.17.45-.3.73-.4a3 3 0 0 1 2.3.21h1.64Zm4.65.76a.24.24 0 0 1-.23.14.42.42 0 0 1-.2-.07 3.59 3.59 0 0 0-.67-.3 1.8 1.8 0 0 0-1.03 0c-.14.05-.27.11-.37.2a.87.87 0 0 0-.23.27.75.75 0 0 0-.08.35c0 .15.04.28.13.39.1.1.21.19.36.27.15.07.32.14.5.2a13.63 13.63 0 0 1 1.16.4c.2.08.36.18.5.3a1.33 1.33 0 0 1 .5 1.07 2 2 0 0 1-.15.78c-.1.24-.25.44-.45.62-.2.17-.43.3-.72.4a3.1 3.1 0 0 1-2.14-.05 2.97 2.97 0 0 1-.87-.53l.25-.41c.04-.05.07-.1.12-.12a.3.3 0 0 1 .17-.04.4.4 0 0 1 .22.08l.3.19a1.91 1.91 0 0 0 1.03.27c.2 0 .38-.03.54-.08.16-.06.29-.13.4-.22a.96.96 0 0 0 .3-.7c0-.17-.05-.31-.14-.42-.09-.11-.2-.2-.36-.28a2.6 2.6 0 0 0-.5-.2l-.59-.19c-.2-.06-.39-.14-.58-.22a2.14 2.14 0 0 1-.5-.3 1.45 1.45 0 0 1-.36-.46c-.1-.19-.14-.41-.14-.67a1.6 1.6 0 0 1 .57-1.23c.18-.16.4-.3.68-.39.26-.1.57-.14.91-.14a2.84 2.84 0 0 1 1.9.7l-.23.4Z"
    />
    <defs>
      <path d="M0 0h85v38H0z"/>
    </defs>
  </svg>
);

export const LogoAnomalyIcon = () => (
  <svg
    viewBox="0 0 85 38"
    fill="currentColor"
  >
    <path
      d="M11.118 10.476c.36.28.801.433 1.257.436h.052c.48-.007.961-.192 1.25-.444 1.509-1.279 5.88-5.287 5.88-5.287 1.168-1.087-2.093-2.174-7.13-2.181h-.06c-5.036.007-8.298 1.094-7.13 2.181 0 0 4.372 4.008 5.88 5.295zm2.559 2.166c-.359.283-.801.439-1.258.444h-.044a2.071 2.071 0 0 1-1.257-.444C10.082 11.755 6.384 8.42 5 7.148v1.93c0 .215.081.496.222.629l.07.064c1.045.955 4.546 4.154 5.825 5.245.358.283.8.438 1.257.444h.044c.489-.015.962-.2 1.258-.444 1.309-1.11 4.948-4.444 5.887-5.31.148-.132.222-.413.222-.628v-1.93a455.127 455.127 0 0 1-6.11 5.494zm-1.258 4.984a2.071 2.071 0 0 0 1.258-.436c2.053-1.815 4.09-3.65 6.11-5.502v1.938c0 .207-.075.488-.223.621-.94.873-4.578 4.2-5.887 5.31-.296.25-.77.436-1.258.443h-.044a2.071 2.071 0 0 1-1.257-.436c-1.204-1.027-4.376-3.928-5.616-5.062l-.28-.255c-.14-.133-.221-.414-.221-.621v-1.938c1.383 1.265 5.081 4.607 6.117 5.495.358.282.8.438 1.257.443h.044zM40 5l-5.84 14.46h-2.43L25.89 5h2.16c.233 0 .423.057.57.17.146.113.256.26.33.44l3.41 8.82c.113.287.22.603.32.95.106.34.206.697.3 1.07.08-.373.166-.73.26-1.07a8.84 8.84 0 0 1 .31-.95l3.39-8.82a.959.959 0 0 1 .31-.42.906.906 0 0 1 .58-.19H40zm17.176 0v14.46h-2.37v-9.34c0-.373.02-.777.06-1.21l-4.37 8.21c-.206.393-.523.59-.95.59h-.38c-.426 0-.743-.197-.95-.59l-4.42-8.24c.02.22.037.437.05.65.014.213.02.41.02.59v9.34h-2.37V5h2.03c.12 0 .224.003.31.01a.778.778 0 0 1 .23.05c.074.027.137.07.19.13.06.06.117.14.17.24l4.33 8.03c.114.213.217.433.31.66.1.227.197.46.29.7.094-.247.19-.483.29-.71.1-.233.207-.457.32-.67l4.27-8.01c.054-.1.11-.18.17-.24a.57.57 0 0 1 .19-.13.903.903 0 0 1 .24-.05c.087-.007.19-.01.31-.01h2.03zm8.887 13.73c.68 0 1.286-.117 1.82-.35.54-.24.996-.57 1.37-.99a4.28 4.28 0 0 0 .85-1.48c.2-.573.3-1.19.3-1.85V5.31h1.02v8.75c0 .78-.124 1.51-.37 2.19a5.248 5.248 0 0 1-1.07 1.77c-.46.5-1.024.893-1.69 1.18-.66.287-1.404.43-2.23.43-.827 0-1.574-.143-2.24-.43a5.012 5.012 0 0 1-1.69-1.18 5.33 5.33 0 0 1-1.06-1.77 6.373 6.373 0 0 1-.37-2.19V5.31h1.03v8.74c0 .66.096 1.277.29 1.85.2.567.483 1.06.85 1.48.373.42.826.75 1.36.99.54.24 1.15.36 1.83.36zm10.38.73h-1.03V5.31h1.03v14.15zM4.242 35v-5.166l-.672-.078a.595.595 0 0 1-.21-.09.23.23 0 0 1-.078-.186v-.438h.96v-.588c0-.348.048-.656.144-.924.1-.272.24-.5.42-.684a1.79 1.79 0 0 1 .66-.426c.256-.096.544-.144.864-.144.272 0 .522.04.75.12l-.024.534c-.008.096-.062.148-.162.156a4.947 4.947 0 0 1-.39.012c-.184 0-.352.024-.504.072a.949.949 0 0 0-.384.234c-.108.108-.192.25-.252.426a2.184 2.184 0 0 0-.084.654v.558h1.752v.774H5.316V35H4.242zM10.416 28.826a3.1 3.1 0 0 1 1.2.222c.356.148.66.358.912.63s.444.602.576.99c.136.384.204.814.204 1.29 0 .48-.068.912-.204 1.296a2.735 2.735 0 0 1-.576.984 2.572 2.572 0 0 1-.912.63 3.175 3.175 0 0 1-1.2.216c-.448 0-.852-.072-1.212-.216a2.572 2.572 0 0 1-.912-.63 2.805 2.805 0 0 1-.582-.984 3.972 3.972 0 0 1-.198-1.296c0-.476.066-.906.198-1.29.136-.388.33-.718.582-.99.252-.272.556-.482.912-.63.36-.148.764-.222 1.212-.222zm0 5.424c.6 0 1.048-.2 1.344-.6.296-.404.444-.966.444-1.686 0-.724-.148-1.288-.444-1.692-.296-.404-.744-.606-1.344-.606-.304 0-.57.052-.798.156a1.507 1.507 0 0 0-.564.45c-.148.196-.26.438-.336.726a3.941 3.941 0 0 0-.108.966c0 .72.148 1.282.444 1.686.3.4.754.6 1.362.6zM15.677 30.14c.192-.416.428-.74.708-.972.28-.236.622-.354 1.026-.354.128 0 .25.014.366.042.12.028.226.072.318.132l-.078.798c-.024.1-.084.15-.18.15-.056 0-.138-.012-.246-.036a1.694 1.694 0 0 0-.366-.036c-.192 0-.364.028-.516.084-.148.056-.282.14-.402.252a1.782 1.782 0 0 0-.318.408c-.092.16-.176.344-.252.552V35h-1.074v-6.078h.612c.116 0 .196.022.24.066.044.044.074.12.09.228l.072.924zM26.761 28.922 24.283 35h-.96l-2.478-6.078h.87a.33.33 0 0 1 .33.222l1.542 3.912c.048.148.09.292.126.432.036.14.07.28.102.42.032-.14.066-.28.102-.42.036-.14.08-.284.132-.432l1.56-3.912a.33.33 0 0 1 .12-.156.311.311 0 0 1 .198-.066h.834zM27.74 35v-6.078h.643c.152 0 .246.074.282.222l.078.624c.224-.276.476-.502.756-.678.28-.176.604-.264.972-.264.408 0 .738.114.99.342.256.228.44.536.552.924.088-.22.2-.41.336-.57a1.987 1.987 0 0 1 1.014-.624c.196-.048.394-.072.594-.072.32 0 .604.052.852.156.252.1.464.248.636.444.176.196.31.438.402.726.092.284.138.61.138.978V35H34.91v-3.87c0-.476-.104-.836-.312-1.08-.208-.248-.508-.372-.9-.372-.176 0-.344.032-.504.096-.156.06-.294.15-.414.27-.12.12-.216.272-.288.456-.068.18-.102.39-.102.63V35h-1.074v-3.87c0-.488-.098-.852-.294-1.092-.196-.24-.482-.36-.858-.36-.264 0-.508.072-.732.216a2.38 2.38 0 0 0-.618.576V35H27.74zM40.746 32.372c-.428.02-.788.058-1.08.114-.292.052-.526.12-.702.204a.923.923 0 0 0-.378.294.639.639 0 0 0-.114.366c0 .26.076.446.228.558.156.112.358.168.606.168.304 0 .566-.054.786-.162.224-.112.442-.28.654-.504v-1.038zm-3.396-2.67c.708-.648 1.56-.972 2.556-.972.36 0 .682.06.966.18.284.116.524.28.72.492.196.208.344.458.444.75.104.292.156.612.156.96V35h-.672a.708.708 0 0 1-.324-.06c-.076-.044-.136-.13-.18-.258l-.132-.444c-.156.14-.308.264-.456.372a2.804 2.804 0 0 1-.462.264c-.16.072-.332.126-.516.162-.18.04-.38.06-.6.06-.26 0-.5-.034-.72-.102a1.618 1.618 0 0 1-.57-.318 1.414 1.414 0 0 1-.372-.522 1.852 1.852 0 0 1-.132-.726 1.419 1.419 0 0 1 .33-.906c.12-.14.274-.272.462-.396s.418-.232.69-.324c.276-.092.596-.166.96-.222.364-.06.78-.096 1.248-.108v-.36c0-.412-.088-.716-.264-.912-.176-.2-.43-.3-.762-.3-.24 0-.44.028-.6.084-.156.056-.294.12-.414.192l-.33.186a.631.631 0 0 1-.324.084.439.439 0 0 1-.264-.078.716.716 0 0 1-.174-.192l-.264-.474zM44.974 29.6c.124-.124.254-.238.39-.342a2.395 2.395 0 0 1 .936-.444c.176-.044.368-.066.576-.066.336 0 .634.058.894.174.26.112.476.272.648.48.176.204.308.45.396.738.092.284.138.598.138.942V35H47.47v-3.918c0-.376-.086-.666-.258-.87-.172-.208-.434-.312-.786-.312-.256 0-.496.058-.72.174a2.58 2.58 0 0 0-.636.474V35h-1.482v-6.156h.906c.192 0 .318.09.378.27l.102.486zM53.085 28.748c.456 0 .87.074 1.242.222a2.692 2.692 0 0 1 1.578 1.626c.144.392.216.83.216 1.314 0 .488-.072.928-.216 1.32-.144.392-.35.726-.618 1.002a2.653 2.653 0 0 1-.96.636 3.333 3.333 0 0 1-1.242.222c-.46 0-.878-.074-1.254-.222a2.712 2.712 0 0 1-.966-.636 2.922 2.922 0 0 1-.618-1.002 3.807 3.807 0 0 1-.216-1.32c0-.484.072-.922.216-1.314.148-.392.354-.724.618-.996.268-.272.59-.482.966-.63a3.397 3.397 0 0 1 1.254-.222zm0 5.202c.512 0 .89-.172 1.134-.516.248-.344.372-.848.372-1.512s-.124-1.17-.372-1.518c-.244-.348-.622-.522-1.134-.522-.52 0-.906.176-1.158.528-.248.348-.372.852-.372 1.512s.124 1.164.372 1.512c.252.344.638.516 1.158.516zM57.252 35v-6.156h.906c.192 0 .318.09.378.27l.096.456c.108-.12.22-.23.336-.33a2.017 2.017 0 0 1 1.32-.492c.388 0 .706.106.954.318.252.208.44.486.564.834a1.93 1.93 0 0 1 .834-.882c.172-.092.354-.16.546-.204.196-.044.392-.066.588-.066.34 0 .642.052.906.156.264.104.486.256.666.456.18.2.316.444.408.732.096.288.144.618.144.99V35h-1.482v-3.918c0-.392-.086-.686-.258-.882-.172-.2-.424-.3-.756-.3-.152 0-.294.026-.426.078a1.026 1.026 0 0 0-.342.228 1.019 1.019 0 0 0-.228.366 1.435 1.435 0 0 0-.084.51V35h-1.488v-3.918c0-.412-.084-.712-.252-.9-.164-.188-.406-.282-.726-.282-.216 0-.418.054-.606.162a1.979 1.979 0 0 0-.516.432V35h-1.482zM70.558 32.372c-.428.02-.788.058-1.08.114-.292.052-.526.12-.702.204a.923.923 0 0 0-.378.294.639.639 0 0 0-.114.366c0 .26.076.446.228.558.156.112.358.168.606.168.304 0 .566-.054.786-.162.224-.112.442-.28.654-.504v-1.038zm-3.396-2.67c.708-.648 1.56-.972 2.556-.972.36 0 .682.06.966.18.284.116.524.28.72.492.196.208.344.458.444.75.104.292.156.612.156.96V35h-.672a.708.708 0 0 1-.324-.06c-.076-.044-.136-.13-.18-.258l-.132-.444c-.156.14-.308.264-.456.372a2.804 2.804 0 0 1-.462.264c-.16.072-.332.126-.516.162-.18.04-.38.06-.6.06-.26 0-.5-.034-.72-.102a1.618 1.618 0 0 1-.57-.318 1.414 1.414 0 0 1-.372-.522 1.852 1.852 0 0 1-.132-.726 1.419 1.419 0 0 1 .33-.906c.12-.14.274-.272.462-.396s.418-.232.69-.324c.276-.092.596-.166.96-.222.364-.06.78-.096 1.248-.108v-.36c0-.412-.088-.716-.264-.912-.176-.2-.43-.3-.762-.3-.24 0-.44.028-.6.084-.156.056-.294.12-.414.192l-.33.186a.631.631 0 0 1-.324.084.439.439 0 0 1-.264-.078.716.716 0 0 1-.174-.192l-.264-.474zM74.9 26.084V35h-1.482v-8.916H74.9zM81.969 28.844l-3.354 7.848a.538.538 0 0 1-.174.234c-.068.056-.174.084-.318.084h-1.104l1.152-2.472-2.49-5.694h1.302c.116 0 .206.028.27.084.068.056.118.12.15.192l1.308 3.192c.044.108.08.216.108.324.032.108.062.218.09.33a32.3 32.3 0 0 1 .108-.33c.036-.112.076-.222.12-.33l1.236-3.186a.437.437 0 0 1 .408-.276h1.188z"
    />
  </svg>
);

export const LogoShortIcon = () => (
  <svg
    viewBox="0 0 15 17"
    fill="currentColor"
  >
    <path
      d="M6.11767 7.47586C6.47736 7.75563 6.91931 7.90898 7.37503 7.91213H7.42681C7.90756 7.90474 8.38832 7.71987 8.67677 7.46846C10.1856 6.18921 14.5568 2.18138 14.5568 2.18138C15.7254 1.09438 12.4637 0.00739 7.42681 0H7.36764C2.3308 0.00739 -0.930935 1.09438 0.237669 2.18138C0.237669 2.18138 4.60884 6.18921 6.11767 7.47586ZM8.67677 9.64243C8.31803 9.92483 7.87599 10.0808 7.41941 10.0861H7.37503C6.91845 10.0808 6.47641 9.92483 6.11767 9.64243C5.0822 8.75513 1.38409 5.42018 0.000989555 4.14832V6.07829C0.000989555 6.29273 0.0823481 6.57372 0.222877 6.70682L0.293316 6.7712L0.293344 6.77122C1.33784 7.72579 4.83903 10.9255 6.11767 12.0161C6.47641 12.2985 6.91845 12.4545 7.37503 12.4597H7.41941C7.90756 12.4449 8.38092 12.2601 8.67677 12.0161C9.9859 10.9069 13.6249 7.57198 14.5642 6.70682C14.7121 6.57372 14.7861 6.29273 14.7861 6.07829V4.14832C12.7662 5.99804 10.7297 7.82949 8.67677 9.64243ZM7.41941 14.6263C7.87513 14.6232 8.31708 14.4698 8.67677 14.19C10.7298 12.3746 12.7663 10.5407 14.7861 8.68853V10.6259C14.7861 10.8329 14.7121 11.1139 14.5642 11.247C13.6249 12.1196 9.9859 15.4471 8.67677 16.5563C8.38092 16.8077 7.90756 16.9926 7.41941 17H7.37503C6.91931 16.9968 6.47736 16.8435 6.11767 16.5637C4.91427 15.5373 1.74219 12.6364 0.502294 11.5025C0.393358 11.4029 0.299337 11.3169 0.222877 11.247C0.0823481 11.1139 0.000989555 10.8329 0.000989555 10.6259V8.68853C1.38409 9.95303 5.0822 13.2953 6.11767 14.1827C6.47641 14.4651 6.91845 14.6211 7.37503 14.6263H7.41941Z"
    />
  </svg>
);

export const SettingsIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M19.14 12.94c.04-.3.06-.61.06-.94 0-.32-.02-.64-.07-.94l2.03-1.58c.18-.14.23-.41.12-.61l-1.92-3.32c-.12-.22-.37-.29-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54c-.04-.24-.24-.41-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.05.3-.09.63-.09.94s.02.64.07.94l-2.03 1.58c-.18.14-.23.41-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z"
    ></path>
  </svg>
);

export const CloseIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M19 6.41 17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"
    ></path>
  </svg>
);

export const RestartIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M12 5V2L8 6l4 4V7c3.31 0 6 2.69 6 6 0 2.97-2.17 5.43-5 5.91v2.02c3.95-.49 7-3.85 7-7.93 0-4.42-3.58-8-8-8zm-6 8c0-1.65.67-3.15 1.76-4.24L6.34 7.34C4.9 8.79 4 10.79 4 13c0 4.08 3.05 7.44 7 7.93v-2.02c-2.83-.48-5-2.94-5-5.91z"
    ></path>
  </svg>
);

export const InfoIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z"></path>
  </svg>
);

export const WarningIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path d="M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z"></path>
  </svg>
);

export const ErrorIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"></path>
  </svg>
);

export const SuccessIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"
    ></path>
  </svg>
);

export const RefreshIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M12 6v3l4-4-4-4v3c-4.42 0-8 3.58-8 8 0 1.57.46 3.03 1.24 4.26L6.7 14.8c-.45-.83-.7-1.79-.7-2.8 0-3.31 2.69-6 6-6zm6.76 1.74L17.3 9.2c.44.84.7 1.79.7 2.8 0 3.31-2.69 6-6 6v-3l-4 4 4 4v-3c4.42 0 8-3.58 8-8 0-1.57-.46-3.03-1.24-4.26z"
    ></path>
  </svg>
);

export const ArrowDownIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path d="M7.41 8.59 12 13.17l4.59-4.58L18 10l-6 6-6-6 1.41-1.41z"></path>
  </svg>
);

export const ArrowDropDownIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path d="m7 10 5 5 5-5z"></path>
  </svg>
);

export const ClockIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M11.99 2C6.47 2 2 6.48 2 12s4.47 10 9.99 10C17.52 22 22 17.52 22 12S17.52 2 11.99 2zM12 20c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8z"
    ></path>
    <path d="M12.5 7H11v6l5.25 3.15.75-1.23-4.5-2.67z"></path>
  </svg>
);

export const CalendarIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M20 3h-1V1h-2v2H7V1H5v2H4c-1.1 0-2 .9-2 2v16c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 18H4V8h16v13z"
    ></path>
  </svg>
);

export const AlarmIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="m22 5.72-4.6-3.86-1.29 1.53 4.6 3.86L22 5.72zM7.88 3.39 6.6 1.86 2 5.71l1.29 1.53 4.59-3.85zM12.5 8H11v6l4.75 2.85.75-1.23-4-2.37V8zM12 4c-4.97 0-9 4.03-9 9s4.02 9 9 9c4.97 0 9-4.03 9-9s-4.03-9-9-9zm0 16c-3.87 0-7-3.13-7-7s3.13-7 7-7 7 3.13 7 7-3.13 7-7 7z"
    ></path>
  </svg>
);

export const KeyboardIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M20 5H4c-1.1 0-1.99.9-1.99 2L2 17c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm-9 3h2v2h-2V8zm0 3h2v2h-2v-2zM8 8h2v2H8V8zm0 3h2v2H8v-2zm-1 2H5v-2h2v2zm0-3H5V8h2v2zm9 7H8v-2h8v2zm0-4h-2v-2h2v2zm0-3h-2V8h2v2zm3 3h-2v-2h2v2zm0-3h-2V8h2v2z"
    ></path>
  </svg>
);

export const PlayIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path d="M8 5v14l11-7z"></path>
  </svg>
);

export const PlayCircleOutlineIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="m10 16.5 6-4.5-6-4.5v9zM12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8z"
    ></path>
  </svg>
);

export const ChartIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path d="m3.5 18.49 6-6.01 4 4L22 6.92l-1.41-1.41-7.09 7.97-4-4L2 16.99z"></path>
  </svg>
);

export const TableIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M10 10.02h5V21h-5zM17 21h3c1.1 0 2-.9 2-2v-9h-5v11zm3-18H5c-1.1 0-2 .9-2 2v3h19V5c0-1.1-.9-2-2-2zM3 19c0 1.1.9 2 2 2h3V10H3v9z"
    ></path>
  </svg>
);

export const CodeIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path d="M9.4 16.6 4.8 12l4.6-4.6L8 6l-6 6 6 6 1.4-1.4zm5.2 0 4.6-4.6-4.6-4.6L16 6l6 6-6 6-1.4-1.4z"></path>
  </svg>
);

export const DeleteIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path d="M6 19c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V7H6v12zM19 4h-3.5l-1-1h-5l-1 1H5v2h14V4z"></path>
  </svg>
);

export const PlusIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"></path>
  </svg>
);

export const MinusIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path d="M19 13H5v-2h14v2z"></path>
  </svg>
);

export const DoneIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M8.9999 14.7854L18.8928 4.8925C19.0803 4.70497 19.3347 4.59961 19.5999 4.59961C19.8651 4.59961 20.1195 4.70497 20.307 4.8925L21.707 6.2925C22.0975 6.68303 22.0975 7.31619 21.707 7.70672L9.70701 19.7067C9.31648 20.0972 8.68332 20.0972 8.2928 19.7067L2.6928 14.1067C2.50526 13.9192 2.3999 13.6648 2.3999 13.3996C2.3999 13.1344 2.50526 12.88 2.6928 12.6925L4.0928 11.2925C4.48332 10.902 5.11648 10.902 5.50701 11.2925L8.9999 14.7854Z"
    />
  </svg>
);

export const VisibilityIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5zM12 17c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5zm0-8c-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3-1.34-3-3-3z"
    ></path>
  </svg>
);

export const VisibilityOffIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M12 7c2.76 0 5 2.24 5 5 0 .65-.13 1.26-.36 1.83l2.92 2.92c1.51-1.26 2.7-2.89 3.43-4.75-1.73-4.39-6-7.5-11-7.5-1.4 0-2.74.25-3.98.7l2.16 2.16C10.74 7.13 11.35 7 12 7zM2 4.27l2.28 2.28.46.46C3.08 8.3 1.78 10.02 1 12c1.73 4.39 6 7.5 11 7.5 1.55 0 3.03-.3 4.38-.84l.42.42L19.73 22 21 20.73 3.27 3 2 4.27zM7.53 9.8l1.55 1.55c-.05.21-.08.43-.08.65 0 1.66 1.34 3 3 3 .22 0 .44-.03.65-.08l1.55 1.55c-.67.33-1.41.53-2.2.53-2.76 0-5-2.24-5-5 0-.79.2-1.53.53-2.2zm4.31-.78 3.15 3.15.02-.16c0-1.66-1.34-3-3-3l-.17.01z"
    ></path>
  </svg>
);

export const Prettify = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M19 9l1.25-2.75L23 5l-2.75-1.25L19 1l-1.25 2.75L15 5l2.75 1.25L19 9zm-7.5.5L9 4 6.5 9.5 1 12l5.5 2.5L9 20l2.5-5.5L17 12l-5.5-2.5zM19 15l-1.25 2.75L15 19l2.75 1.25L19 23l1.25-2.75L23 19l-2.75-1.25L19 15z"
    ></path>
  </svg>
);

export const CopyIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z"
    ></path>
  </svg>
);

export const DragIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path d="M20 9H4v2h16V9zM4 15h16v-2H4v2z"></path>
  </svg>
);

export const TimelineIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M23 8c0 1.1-.9 2-2 2-.18 0-.35-.02-.51-.07l-3.56 3.55c.05.16.07.34.07.52 0 1.1-.9 2-2 2s-2-.9-2-2c0-.18.02-.36.07-.52l-2.55-2.55c-.16.05-.34.07-.52.07s-.36-.02-.52-.07l-4.55 4.56c.05.16.07.33.07.51 0 1.1-.9 2-2 2s-2-.9-2-2 .9-2 2-2c.18 0 .35.02.51.07l4.56-4.55C8.02 9.36 8 9.18 8 9c0-1.1.9-2 2-2s2 .9 2 2c0 .18-.02.36-.07.52l2.55 2.55c.16-.05.34-.07.52-.07s.36.02.52.07l3.55-3.56C19.02 8.35 19 8.18 19 8c0-1.1.9-2 2-2s2 .9 2 2z"
    ></path>
  </svg>
);

export const WikiIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M21 5C19.89 4.65 18.67 4.5 17.5 4.5C15.55 4.5 13.45 4.9 12 6C10.55 4.9 8.45 4.5 6.5 4.5C5.33 4.5 4.11 4.65 3 5C2.25 5.25 1.6 5.55 1 6V20.6C1 20.85 1.25 21.1 1.5 21.1C1.6 21.1 1.65 21.1 1.75 21.05C3.15 20.3 4.85 20 6.5 20C8.2 20 10.65 20.65 12 21.5C13.35 20.65 15.8 20 17.5 20C19.15 20 20.85 20.3 22.25 21.05C22.35 21.1 22.4 21.1 22.5 21.1C22.75 21.1 23 20.85 23 20.6V6C22.4 5.55 21.75 5.25 21 5ZM21 18.5C19.9 18.15 18.7 18 17.5 18C15.8 18 13.35 18.65 12 19.5C10.65 18.65 8.2 18 6.5 18C5.3 18 4.1 18.15 3 18.5V7C4.1 6.65 5.3 6.5 6.5 6.5C8.2 6.5 10.65 7.15 12 8C13.35 7.15 15.8 6.5 17.5 6.5C18.7 6.5 19.9 6.65 21 7V18.5Z"
    />
    <path
      d="M17.5 10.5C18.38 10.5 19.23 10.59 20 10.76V9.24C19.21 9.09 18.36 9 17.5 9C15.8 9 14.26 9.29 13 9.83V11.49C14.13 10.85 15.7 10.5 17.5 10.5ZM13 12.49V14.15C14.13 13.51 15.7 13.16 17.5 13.16C18.38 13.16 19.23 13.25 20 13.42V11.9C19.21 11.75 18.36 11.66 17.5 11.66C15.8 11.66 14.26 11.96 13 12.49ZM17.5 14.33C15.8 14.33 14.26 14.62 13 15.16V16.82C14.13 16.18 15.7 15.83 17.5 15.83C18.38 15.83 19.23 15.92 20 16.09V14.57C19.21 14.41 18.36 14.33 17.5 14.33Z"
    />
    <path
      d="M6.5 10.5C5.62 10.5 4.77 10.59 4 10.76V9.24C4.79 9.09 5.64 9 6.5 9C8.2 9 9.74 9.29 11 9.83V11.49C9.87 10.85 8.3 10.5 6.5 10.5ZM11 12.49V14.15C9.87 13.51 8.3 13.16 6.5 13.16C5.62 13.16 4.77 13.25 4 13.42V11.9C4.79 11.75 5.64 11.66 6.5 11.66C8.2 11.66 9.74 11.96 11 12.49ZM6.5 14.33C8.2 14.33 9.74 14.62 11 15.16V16.82C9.87 16.18 8.3 15.83 6.5 15.83C5.62 15.83 4.77 15.92 4 16.09V14.57C4.79 14.41 5.64 14.33 6.5 14.33Z"
    />
  </svg>

);

export const IssueIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M12 2C6.49 2 2 6.49 2 12s4.49 10 10 10 10-4.49 10-10S17.51 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm3-8c0 1.66-1.34 3-3 3s-3-1.34-3-3 1.34-3 3-3 3 1.34 3 3z"
    ></path>
  </svg>
);

export const QuestionIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M12 2C6.48 2 2 6.48 2 12C2 17.52 6.48 22 12 22C17.52 22 22 17.52 22 12C22 6.48 17.52 2 12 2ZM12 6C9.79 6 8 7.79 8 10H10C10 8.9 10.9 8 12 8C13.1 8 14 8.9 14 10C14 10.8792 13.4202 11.3236 12.7704 11.8217C11.9421 12.4566 11 13.1787 11 15H13C13 13.9046 13.711 13.2833 14.4408 12.6455C15.21 11.9733 16 11.2829 16 10C16 7.79 14.21 6 12 6ZM13 16V18H11V16H13Z"
    />
  </svg>

);

export const StorageIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M4 20h16c1.1 0 2-.9 2-2s-.9-2-2-2H4c-1.1 0-2 .9-2 2s.9 2 2 2zm0-3h2v2H4v-2zM2 6c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2s-.9-2-2-2H4c-1.1 0-2 .9-2 2zm4 1H4V5h2v2zm-2 7h16c1.1 0 2-.9 2-2s-.9-2-2-2H4c-1.1 0-2 .9-2 2s.9 2 2 2zm0-3h2v2H4v-2z"
    ></path>
  </svg>
);

export const MoreIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M12 8c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2zm0 2c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm0 6c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z"
    ></path>
  </svg>
);

export const TuneIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M3 17v2h6v-2H3zM3 5v2h10V5H3zm10 16v-2h8v-2h-8v-2h-2v6h2zM7 9v2H3v2h4v2h2V9H7zm14 4v-2H11v2h10zm-6-4h2V7h4V5h-4V3h-2v6z"
    ></path>
  </svg>
);

export const TipIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M7 20h4c0 1.1-.9 2-2 2s-2-.9-2-2zm-2-1h8v-2H5v2zm11.5-9.5c0 3.82-2.66 5.86-3.77 6.5H5.27c-1.11-.64-3.77-2.68-3.77-6.5C1.5 5.36 4.86 2 9 2s7.5 3.36 7.5 7.5zm4.87-2.13L20 8l1.37.63L22 10l.63-1.37L24 8l-1.37-.63L22 6l-.63 1.37zM19 6l.94-2.06L22 3l-2.06-.94L19 0l-.94 2.06L16 3l2.06.94L19 6z"
    ></path>
  </svg>
);

export const ListIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path d="M3 14h4v-4H3v4zm0 5h4v-4H3v4zM3 9h4V5H3v4zm5 5h13v-4H8v4zm0 5h13v-4H8v4zM8 5v4h13V5H8z"></path>
  </svg>
);

export const StarBorderIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="m22 9.24-7.19-.62L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21 12 17.27 18.18 21l-1.63-7.03L22 9.24zM12 15.4l-3.76 2.27 1-4.28-3.32-2.88 4.38-.38L12 6.1l1.71 4.04 4.38.38-3.32 2.88 1 4.28L12 15.4z"
    ></path>
  </svg>
);

export const StarIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path d="M12 17.27 18.18 21l-1.64-7.03L22 9.24l-7.19-.61L12 2 9.19 8.63 2 9.24l5.46 4.73L5.82 21z"></path>
  </svg>
);

export const MetricIcon = () => (
  <svg
    viewBox="0 0 16 16"
    fill={getCssVariable("color-error")}
  >
    <path
      d="M13.5095 4L8.50952 1H7.50952L2.50952 4L2.01953 4.85999V10.86L2.50952 11.71L7.50952 14.71H8.50952L13.5095 11.71L13.9995 10.86V4.85999L13.5095 4ZM7.50952 13.5601L3.00952 10.86V5.69995L7.50952 8.15002V13.5601ZM3.26953 4.69995L8.00952 1.85999L12.7495 4.69995L8.00952 7.29004L3.26953 4.69995ZM13.0095 10.86L8.50952 13.5601V8.15002L13.0095 5.69995V10.86Z"
    />
  </svg>
);

export const FunctionIcon = () => (
  <svg
    viewBox="0 0 16 16"
    fill={getCssVariable("color-primary")}
  >
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M2 5H4V4H1.5L1 4.5V12.5L1.5 13H4V12H2V5ZM14.5 4H12V5H14V12H12V13H14.5L15 12.5V4.5L14.5 4ZM11.76 6.56995L12 7V9.51001L11.7 9.95996L7.19995 11.96H6.73999L4.23999 10.46L4 10.03V7.53003L4.30005 7.06995L8.80005 5.06995H9.26001L11.76 6.56995ZM5 9.70996L6.5 10.61V9.28003L5 8.38V9.70996ZM5.57996 7.56006L7.03003 8.43005L10.42 6.93005L8.96997 6.06006L5.57996 7.56006ZM7.53003 10.73L11.03 9.17004V7.77002L7.53003 9.31995V10.73Z"
    />
  </svg>
);

export const LabelIcon = () => (
  <svg
    viewBox="0 0 16 16"
    fill={getCssVariable("color-warning")}
  >
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M14 2H8L7 3V6H8V3H14V8H10V9H14L15 8V3L14 2ZM9 6H13V7H9.41L9 6.59V6ZM7 7H2L1 8V13L2 14H8L9 13V8L8 7H7ZM8 13H2V8H8V9V13ZM3 9H7V10H3V9ZM3 11H7V12H3V11ZM9 4H13V5H9V4Z"
    />
  </svg>
);

export const ValueIcon = () => (
  <svg
    viewBox="0 0 16 16"
    fill={getCssVariable("color-primary")}
  >
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M7 3L8 2H14L15 3V8L14 9H10V8H14V3H8V6H7V3ZM9 9V8L8 7H7H2L1 8V13L2 14H8L9 13V9ZM8 8V9V13H2V8H7H8ZM9.41421 7L9 6.58579V6H13V7H9.41421ZM9 4H13V5H9V4ZM7 10H3V11H7V10Z"
    />
  </svg>
);

export const DownloadIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path d="M19 9h-4V3H9v6H5l7 7 7-7zM5 18v2h14v-2H5z"></path>
  </svg>
);

export const ExpandIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M12 5.83 15.17 9l1.41-1.41L12 3 7.41 7.59 8.83 9zm0 12.34L8.83 15l-1.41 1.41L12 21l4.59-4.59L15.17 15z"
    ></path>
  </svg>
);

export const CollapseIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M7.41 18.59 8.83 20 12 16.83 15.17 20l1.41-1.41L12 14zm9.18-13.18L15.17 4 12 7.17 8.83 4 7.41 5.41 12 10z"
    ></path>
  </svg>
);

export const SearchIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14"
    ></path>
  </svg>
);

export const SpinnerIcon = () => (
  <svg viewBox="0 0 24 24">
    <path
      fill="currentColor"
      d="M12,4a8,8,0,0,1,7.89,6.7A1.53,1.53,0,0,0,21.38,12h0a1.5,1.5,0,0,0,1.48-1.75,11,11,0,0,0-21.72,0A1.5,1.5,0,0,0,2.62,12h0a1.53,1.53,0,0,0,1.49-1.3A8,8,0,0,1,12,4Z"
    >
      <animateTransform
        attributeName="transform"
        dur="0.75s"
        repeatCount="indefinite"
        type="rotate"
        values="0 12 12;360 12 12"
      />
    </path>
  </svg>
);

export const CommentIcon = () => (
  <svg
    viewBox="0 0 24 24"
    fill="currentColor"
  >
    <path
      d="M21.99 4c0-1.1-.89-2-1.99-2H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h14l4 4zM18 14H6v-2h12zm0-3H6V9h12zm0-3H6V6h12z"
    ></path>
  </svg>
);
