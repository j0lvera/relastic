# React CRUD Generator for Elastic UI

A command-line tool that generates React CRUD (Create, Read, Update, Delete) components using the Elastic UI component library. 

This tool helps streamline the development process by automatically generating consistent, type-safe React components.

## Features

- Generates complete CRUD functionality
- Uses TypeScript for type safety
- Integrates with Elastic UI components
- Includes:
    - TypeScript interfaces and types
    - API integration with React Query
    - Form components with validation
    - Table components with sorting and pagination
    - List views with search functionality
    - Detail views with related operations
    - Modal components for create/update/delete operations

## Installation

TODO

## Usage

TODO

Field definitions follow the format: `name:type:required`, where:
- `name`: The field name
- `type`: One of `string`, `number`, or `boolean`
- `required`: `true` or `false`

Multiple fields are separated by commas.

### Example

```bash
./relastic User "name:string:true,age:number:false,isActive:boolean:true"
```

This will generate the following files in a `user` directory:
- `User.types.ts` - TypeScript interfaces and types
- `User.constants.ts` - API endpoints and query keys
- `User.api.ts` - API integration with React Query
- `User.form.tsx` - Form component for create/update
- `User.table.tsx` - Table component with sorting and pagination
- `User.list.tsx` - List view with search and modals
- `User.detail.tsx` - Detail view with related operations

## Generated Components

### Types
Defines TypeScript interfaces and types for your entity, including:
- Base interface
- Create/Update/Delete types
- Table props
- Form props
- Pagination types

### API Integration
Includes React Query hooks for:
- Fetching list data
- Fetching single entity
- Creating new entities
- Updating existing entities
- Deleting entities

### Form Component
- Field validation
- Error handling
- Type-safe form submission
- Integration with Elastic UI form components

### Table Component
- Sortable columns
- Pagination
- Row selection
- Edit/Delete actions
- Type-safe data handling

### List Component
- Search functionality
- Create/Edit/Delete modals
- Pagination handling
- Integration with React Router

### Detail Component
- Single entity display
- Related operations
- Breadcrumb navigation
- Modal handling for operations

## Dependencies

The generated components expect the following dependencies in your React project:
- `@elastic/eui`
- `@tanstack/react-query`
- `@tanstack/react-router`
- `react-hook-form`

## License

This project is licensed under the MIT License.