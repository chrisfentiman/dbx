# Geospatial Conventions

## CRS (Coordinate Reference Systems)
- ALWAYS set CRS explicitly on every GeoDataFrame
- Use EPSG:4326 (WGS84) for storage and data interchange
- Project to local UTM or appropriate metric CRS for distance/area calculations
- Never calculate distance or area in geographic (degree-based) CRS
- Use `gdf.estimate_utm_crs()` when a local projected CRS is needed dynamically
- Document CRS assumptions in function docstrings

## Geometry Operations
- Use Shapely 2.x vectorized functions (not `.apply()` on geometry columns)
- Validate geometries with `.is_valid` before topology operations
- Use `.make_valid()` to fix invalid geometries
- Handle null/empty geometries explicitly before joins and operations
- Prefer `unary_union` over iterative union

## Spatial Joins and Queries
- Use `gpd.sjoin()` with explicit `predicate=` parameter
- Use `gpd.sjoin_nearest()` for proximity queries
- Filter with bounding box BEFORE expensive spatial operations
- Use spatial indexes (`.sindex`) for custom spatial queries

## Data Formats
- GeoParquet for large datasets (fast, columnar, preserves CRS)
- GeoJSON for small datasets and API responses
- GeoPackage for multi-layer vector data
- Cloud-Optimized GeoTIFF (COG) for raster data
- Never use Shapefile for new work (limited field names, no null support, multi-file mess)

## PostGIS
- Always use parameterized queries (never string interpolation for SQL)
- Use `ST_Intersects` with bounding box for spatial filtering
- Create spatial indexes: `CREATE INDEX ON table USING GIST (geom)`
- Use `geography` type for global distance queries, `geometry` for local analysis

## Performance
- Filter spatially before loading full datasets
- Use `bbox` parameter in `gpd.read_file()` and `gpd.read_parquet()`
- Chunk large datasets rather than loading everything into memory
- Prefer GeoParquet over GeoJSON for datasets > 10MB
- Use DuckDB spatial extension for analytical queries on large files

## Visualization
- Folium for interactive web maps
- Pydeck for 3D and large-scale visualization
- Include basemap context (use contextily for static maps)
- Always include a legend and CRS info in map outputs
