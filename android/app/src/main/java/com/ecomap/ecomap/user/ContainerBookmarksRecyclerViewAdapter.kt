package com.ecomap.ecomap.user

import android.content.Context
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageButton
import android.widget.TextView
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.ecomap.ecomap.R
import com.ecomap.ecomap.domain.Container
import com.ecomap.ecomap.map.ContainerCategoriesRecyclerViewAdapter
import com.ecomap.ecomap.map.ContainerCategoryRecyclerViewData

/**
 * Represents the structure of an item in the container bookmark recycler view.
 * @param containers Containers associated with the bookmark.
 */
data class ContainerBookmarkRecyclerViewData(val containers: ArrayList<Container>)

/**
 * Recycler view adapter for the container bookmarks.
 */
class ContainerBookmarksRecyclerViewAdapter(
    private val context: Context,
    private val dataSet: ArrayList<ContainerBookmarkRecyclerViewData>
) :
    RecyclerView.Adapter<ContainerBookmarksRecyclerViewAdapter.ViewHolder>() {
    var onButtonContainerBookmarkClicked: ((position: Int) -> Unit)? = null

    /**
     * Defines the views in the adapter.
     */
    class ViewHolder(view: View) : RecyclerView.ViewHolder(view) {
        val textViewMunicipalityName: TextView = view.findViewById(R.id.text_view_municipality_name)
        val textViewWayName: TextView = view.findViewById(R.id.text_view_way_name)
        val buttonContainerBookmark: ImageButton =
            view.findViewById(R.id.button_container_bookmark)
        val recyclerContainerCategories: RecyclerView =
            view.findViewById(R.id.recycler_container_categories)
    }

    override fun onCreateViewHolder(viewGroup: ViewGroup, viewType: Int): ViewHolder {
        // Create a new view, which defines the UI of the list item.
        val view = LayoutInflater.from(viewGroup.context)
            .inflate(R.layout.frame_layout_container_bookmark, viewGroup, false)

        return ViewHolder(view)
    }

    override fun onBindViewHolder(viewHolder: ViewHolder, position: Int) {
        val data = dataSet[position]

        // Set location text data.
        if (data.containers.isNotEmpty()) {
            val container = data.containers[0]

            viewHolder.textViewMunicipalityName.text = container.geoJSON.properties.municipalityName
            viewHolder.textViewWayName.text = container.geoJSON.properties.getWayName(context)
        }

        // Set button functions.
        viewHolder.buttonContainerBookmark.setOnClickListener {
            onButtonContainerBookmarkClicked?.invoke(position)
        }

        // Populate the container category recycler view.
        val containerCategoriesDataSet =
            ArrayList<ContainerCategoryRecyclerViewData>(data.containers.size)
        for (container in data.containers) {
            val categoryData =
                ContainerCategoryRecyclerViewData(container.category.getIconResource())
            if (containerCategoriesDataSet.contains(categoryData)) {
                // The category already exists in the current data set.
                continue
            }

            containerCategoriesDataSet.add(categoryData)
        }

        viewHolder.recyclerContainerCategories.layoutManager =
            LinearLayoutManager(context, RecyclerView.HORIZONTAL, false)
        viewHolder.recyclerContainerCategories.adapter =
            ContainerCategoriesRecyclerViewAdapter(containerCategoriesDataSet.toTypedArray())
    }

    override fun getItemCount(): Int {
        return dataSet.size
    }
}
